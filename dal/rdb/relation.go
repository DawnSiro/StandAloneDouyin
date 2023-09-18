package rdb

import (
	"douyin/dal/model"
	"douyin/pkg/constant"
	"douyin/pkg/global"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis"
	"strconv"
)

// RelationZSet 集合
// RelationID 做为 Score
// ToUserID or UserID 做为 Member
type RelationZSet struct {
	ID     uint64
	Member uint64
}

type FollowUserZSet struct {
	ID     uint64
	Member uint64
}

// Follow 关注操作，使用 relation 表的主键做为排序依据
// 1. 关注者的关注 set 集合增加
// 2. 关注者的 ZSet 关注用户信息集合增加
// 3. 关注者的关注数+1
// 4. 被关注者的粉丝数+1
// 5. 被关注者的 ZSet 粉丝信息集合增加
func Follow(user *model.FanUserRedisData, toUser *model.FollowUserData) error {
	logTag := "dal.rdb.relation.Follow err:"
	// 生成参数
	followIDKey := constant.FollowIDRedisSetPrefix + strconv.FormatUint(user.UID, 10)
	followKey := constant.FollowRedisZSetPrefix + strconv.FormatUint(user.UID, 10)
	followInfoCountKey := constant.UserInfoCountHashRedisPrefix + strconv.FormatUint(user.UID, 10)
	fanKey := constant.FanRedisZSetPrefix + strconv.FormatUint(toUser.UID, 10)
	fanInfoCountKey := constant.UserInfoCountHashRedisPrefix + strconv.FormatUint(toUser.UID, 10)

	userJSON, err := json.Marshal(user)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return err
	}

	// ZSet 中的用户信息
	toUserRedisData := &model.FollowUserRedisData{
		UID:      toUser.UID,
		Username: toUser.Username,
		Avatar:   toUser.Avatar,
	}

	toUserJSON, err := json.Marshal(toUserRedisData)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return err
	}

	keys := []string{followIDKey, followKey, followInfoCountKey, fanKey, fanInfoCountKey,
		constant.FollowCountRedisFiled, constant.FanCountRedisFiled}
	args := []interface{}{float64(toUser.CreatedTime.UnixMilli()), toUser.UID, userJSON, toUserJSON}
	// 执行 Redis 服务端缓存的 Lua 脚本，Lua 脚本放在 constant 包下的 redis.go 中
	err = global.UserRC.EvalSha(global.FollowLuaScriptHash, keys, args...).Err()
	if err != nil {
		hlog.Error(logTag, err.Error())
		return err
	}
	return nil
}

// CancelFollow 取消关注操作，此时无需传入做为排序依据的 relationID
func CancelFollow(user *model.FanUserRedisData, toUser *model.FollowUserData) error {
	logTag := "dal.rdb.relation.CancelFollow err:"
	followIDKey := constant.FollowIDRedisSetPrefix + strconv.FormatUint(user.UID, 10)
	followKey := constant.FollowRedisZSetPrefix + strconv.FormatUint(user.UID, 10)
	followInfoCountKey := constant.UserInfoCountHashRedisPrefix + strconv.FormatUint(user.UID, 10)
	fanKey := constant.FanRedisZSetPrefix + strconv.FormatUint(toUser.UID, 10)
	fanInfoCountKey := constant.UserInfoCountHashRedisPrefix + strconv.FormatUint(toUser.UID, 10)

	userJSON, err := json.Marshal(user)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return err
	}

	// ZSet 中的用户信息
	toUserRedisData := &model.FollowUserRedisData{
		UID:      toUser.UID,
		Username: toUser.Username,
		Avatar:   toUser.Avatar,
	}

	toUserJSON, err := json.Marshal(toUserRedisData)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return err
	}

	keys := []string{followIDKey, followKey, followInfoCountKey, fanKey, fanInfoCountKey,
		constant.FollowCountRedisFiled, constant.FanCountRedisFiled}
	args := []interface{}{float64(toUser.CreatedTime.UnixMilli()), toUser.UID, userJSON, toUserJSON}
	// 执行 Redis 服务端缓存的 Lua 脚本
	err = global.UserRC.EvalSha(global.CancelFollowLuaScriptHash, keys, args...).Err()
	if err != nil {
		hlog.Error(logTag, err.Error())
		return err
	}
	return nil
}

// SetFollowUserZSet 设置关注用户ID ZSet
// 一百个以上的 value 容易导致设置缓慢，故每次都循环设置一百个
func SetFollowUserZSet(userID uint64, followList []*model.FollowUserData) error {
	logTag := "dal.rdb.relation.SetFollowUserZSet err:"
	key := constant.FollowRedisZSetPrefix + strconv.FormatUint(userID, 10)
	zList := make([]redis.Z, len(followList))
	for i := 0; i < len(followList); i++ {
		zList[i].Score = float64(followList[i].CreatedTime.UnixMilli())
		followUserJSON, err := json.Marshal(model.FollowUserRedisData{
			UID:      followList[i].UID,
			Username: followList[i].Username,
			Avatar:   followList[i].Avatar,
		})
		if err != nil {
			hlog.Error(logTag, err.Error())
			return err
		}
		zList[i].Member = followUserJSON
	}
	// 批量地添加值，每次一百个
	var err error
	for i := 0; i < len(followList); i += 100 {
		if len(followList)-i > 100 {
			err = global.UserRC.ZAdd(key, zList[i:i+100]...).Err()
		} else {
			err = global.UserRC.ZAdd(key, zList[i:len(followList)]...).Err()
		}
		if err != nil {
			hlog.Error(logTag, err.Error())
			// 删除整个 ZSet
			global.UserRC.Del(key)
			return err
		}
	}
	return nil
}

// GetFollowUserZSet 获取信息
func GetFollowUserZSet(userID uint64) ([]*model.FollowUserRedisData, error) {
	key := constant.FollowRedisZSetPrefix + strconv.FormatUint(userID, 10)
	// 先判断缓存是否存在
	exists, _ := global.UserRC.Exists(key).Result()
	if exists == 0 {
		return nil, redis.Nil
	}

	zSet, err := global.UserRC.ZRevRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	// 转换结果
	result := make([]*model.FollowUserRedisData, len(zSet))
	for i, v := range zSet {
		err = json.Unmarshal([]byte(v), &result[i])
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// SetFanUserZSet 设置关注用户ID ZSet
// 一百个以上的 value 容易导致设置缓慢，故每次都循环设置一百个
func SetFanUserZSet(userID uint64, followList []*model.FanUserData) error {
	logTag := "dal.rdb.relation.SetFanUserZSet err:"
	key := constant.FanRedisZSetPrefix + strconv.FormatUint(userID, 10)
	zList := make([]redis.Z, len(followList))
	for i := 0; i < len(followList); i++ {
		zList[i].Score = float64(followList[i].CreatedTime.UnixMilli())
		followUserJSON, err := json.Marshal(model.FollowUserRedisData{
			UID:      followList[i].UID,
			Username: followList[i].Username,
			Avatar:   followList[i].Avatar,
		})
		if err != nil {
			hlog.Error(logTag, err.Error())
			return err
		}
		zList[i].Member = followUserJSON
	}
	// 批量地添加值，每次一百个
	var err error
	for i := 0; i < len(followList); i += 100 {
		if len(followList)-i > 100 {
			err = global.UserRC.ZAdd(key, zList[i:i+100]...).Err()
		} else {
			err = global.UserRC.ZAdd(key, zList[i:len(followList)]...).Err()
		}
		if err != nil {
			hlog.Error(logTag, err.Error())
			// 删除整个 ZSet
			global.UserRC.Del(key)
			return err
		}
	}
	return nil
}

// GetFanUserZSet 获取信息
func GetFanUserZSet(userID uint64) ([]*model.FanUserRedisData, error) {
	logTag := "dal.rdb.relation.GetFanUserZSet err:"
	key := constant.FanRedisZSetPrefix + strconv.FormatUint(userID, 10)
	zSet, err := global.UserRC.ZRevRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	// 转换结果
	result := make([]*model.FanUserRedisData, len(zSet))
	for i, v := range zSet {
		err = json.Unmarshal([]byte(v), &result[i])
		if err != nil {
			hlog.Error(logTag, err.Error())
			return nil, err
		}
	}
	return result, nil
}

func GetFriendIDZSet(userID uint64) ([]uint64, error) {
	followKey := constant.FollowRedisZSetPrefix + strconv.FormatUint(userID, 10)
	fanKey := constant.FanRedisZSetPrefix + strconv.FormatUint(userID, 10)

	// 使用 ZInterStore 方法计算交集，并将结果存储在一个临时的键中
	// ZInterStore temp:userID 2 followRedisZSetPrefix:userID fanRedisZSetPrefix:userID WEIGHTS 1 1 AGGREGATE MAX
	tempKey := "temp:" + strconv.FormatUint(userID, 10)
	zInterStoreArgs := redis.ZStore{
		// 权重，这里两个集合权重对等
		Weights: []float64{1, 1},
		// 指定聚合方式，这里使用 MAX，也可以使用 MIN 或 SUM
		Aggregate: "MAX",
	}
	err := global.UserRC.ZInterStore(tempKey, zInterStoreArgs, followKey, fanKey).Err()
	if err != nil {
		return nil, err
	}

	// 获取临时键中的成员
	zSet, err := global.UserRC.ZRevRange(tempKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// 删除临时键
	global.UserRC.Del(tempKey)

	// 将结果转换为 uint64 切片
	result := make([]uint64, len(zSet))
	for i, v := range zSet {
		result[i], _ = strconv.ParseUint(v, 10, 64)
	}

	return result, nil
}

func LoadFollowZSet(userID uint64, r []*RelationZSet) error {
	key := constant.FollowRedisZSetPrefix + strconv.FormatUint(userID, 10)
	// 预先申请够足够的内存大小，避免调用 append 函数带来的性能损耗
	redisZ := make([]redis.Z, len(r))
	for i := 0; i < len(r); i++ {
		redisZ[i].Score = float64(r[i].ID)
		redisZ[i].Member = r[i].Member
	}
	return global.UserRC.ZAdd(key, redisZ...).Err()
}

func LoadFanZSet(userID uint64, r []*RelationZSet) error {
	key := constant.FanRedisZSetPrefix + strconv.FormatUint(userID, 10)
	// 预先申请够足够的内存大小，避免调用 append 函数带来的性能损耗
	redisZ := make([]redis.Z, len(r))
	for i := 0; i < len(r); i++ {
		redisZ[i].Score = float64(r[i].ID)
		redisZ[i].Member = r[i].Member
	}
	return global.UserRC.ZAdd(key, redisZ...).Err()
}

// IsFollow 判断 user 是否关注了 toUser
func IsFollow(userID, toUserID uint64) (bool, error) {
	key := constant.FollowRedisZSetPrefix + strconv.FormatUint(userID, 10)
	// 先判断键是否存在
	exist, _ := global.UserRC.Exists(key).Result()
	if exist > 0 {
		// 键存在的话在进行查询
		return global.UserRC.SIsMember(key, strconv.FormatUint(toUserID, 10)).Result()
	}
	// 键不存在返回 redis.Nil
	return false, redis.Nil
}

// SetFollowUserIDSet 设置关注列表用户ID
func SetFollowUserIDSet(userID uint64, followIDSet []uint64) error {
	key := constant.FollowIDRedisSetPrefix + strconv.FormatUint(userID, 10)

	// 将 []uint64 转换为 []string
	strFollowIDSet := make([]string, len(followIDSet))
	for i, id := range followIDSet {
		strFollowIDSet[i] = strconv.FormatUint(id, 10)
	}

	return global.UserRC.SAdd(key, strFollowIDSet).Err()
}

// GetFollowUserIDSet 获取关注列表用户ID
func GetFollowUserIDSet(userID uint64) (map[uint64]struct{}, error) {
	key := constant.FollowIDRedisSetPrefix + strconv.FormatUint(userID, 10)
	memberString, err := global.UserRC.SMembersMap(key).Result()
	if err != nil {
		return nil, err
	}
	followIDSet, err := convertStringMapToUint64Map(memberString)
	if err != nil {
		return nil, err
	}
	return followIDSet, nil
}

func SetFanUserIDSet(userID uint64, fanIDSet []uint64) error {
	key := constant.FanIDRedisSetPrefix + strconv.FormatUint(userID, 10)
	return global.UserRC.SAdd(key, fanIDSet).Err()
}

func GetFanUserIDSet(userID uint64) (map[uint64]struct{}, error) {
	key := constant.FanIDRedisSetPrefix + strconv.FormatUint(userID, 10)
	memberString, err := global.UserRC.SMembersMap(key).Result()
	if err != nil {
		return nil, err
	}
	fanIDSet, err := convertStringMapToUint64Map(memberString)
	if err != nil {
		return nil, err
	}
	return fanIDSet, nil
}

func convertStringMapToUint64Map(inputMap map[string]struct{}) (map[uint64]struct{}, error) {
	// 提前为 outputMap 分配内存空间
	outputMap := make(map[uint64]struct{}, len(inputMap))
	for key := range inputMap {
		num, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			return nil, err
		}
		outputMap[num] = struct{}{}
	}
	return outputMap, nil
}
