package rdb

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"douyin/pkg/constant"
	"douyin/pkg/global"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis"
)

type CommentInfo struct {
	ID          uint64
	VideoID     uint64
	UserID      uint64
	Content     string
	CreatedTime float64
}

type CommentIDZSet struct {
	CommentID  uint64
	CreateTime int64
}

type CommentRedisZSetData struct {
	CID      uint64 `gorm:"column:cid"`
	Content  string
	UID      uint64
	Username string
	Avatar   string
}

func AddComment(videoID uint64, comment CommentInfo) error {
	jsonString, err := json.Marshal(comment)
	if err != nil {
		hlog.Error("dal.rdb.comment.AddComment err:", err.Error())
		return err
	}
	infoKey := constant.CommentInfoRedisPrefix + strconv.FormatUint(comment.ID, 10)
	zSetKey := constant.CommentRedisZSetPrefix + strconv.FormatUint(videoID, 10)
	keys := []string{infoKey, zSetKey}
	args := []interface{}{jsonString,
		int64(constant.CommentInfoExpiration + time.Duration(rand.Intn(200))*time.Minute),
		comment.CreatedTime, comment.ID}
	// 执行 Redis 服务端缓存的 Lua 脚本
	err = global.UserRC.EvalSha(global.CommentLuaScriptHash, keys, args...).Err()
	if err != nil {
		hlog.Error("dal.rdb.comment.AddComment err:", err.Error())
		return err
	}
	return nil
}

func DeleteComment(videoID, commentID uint64) error {
	infoKey := constant.CommentInfoRedisPrefix + strconv.FormatUint(commentID, 10)
	zSetKey := constant.CommentRedisZSetPrefix + strconv.FormatUint(videoID, 10)
	keys := []string{infoKey, zSetKey}
	args := []interface{}{commentID}
	// 执行 Redis 服务端缓存的 Lua 脚本
	err := global.UserRC.EvalSha(global.DeleteCommentLuaScriptHash, keys, args...).Err()
	if err != nil {
		hlog.Error("dal.rdb.comment.DeleteComment err:", err.Error())
		return err
	}
	return nil
}

// IsCommentCreatedByMyself 获取评论用户ID
func IsCommentCreatedByMyself(userID, commentID uint64) (bool, error) {
	key := constant.CommentInfoRedisPrefix + strconv.FormatUint(commentID, 10)
	// 这里查询的是一个 string ，查不到就是查不到
	result, err := global.CommentRC.Get(key).Result()
	if err != nil {
		hlog.Error("dal.rdb.comment.IsCommentCreatedByMyself err:", err.Error())
		return false, err
	}
	var commentInfo CommentInfo
	err = json.Unmarshal([]byte(result), &commentInfo)
	if err != nil {
		hlog.Error("dal.rdb.comment.IsCommentCreatedByMyself err:", err.Error())
		return false, err
	}
	// 解析完 JSON 判断是否相等，不相等就不是自己发的
	return userID == commentInfo.UserID, nil
}

func SetCommentIDByVideoID(videoID uint64, cIDZ []*CommentIDZSet) error {
	key := constant.CommentRedisZSetPrefix + strconv.FormatUint(videoID, 10)
	global.CommentRC.ZAdd(key)
	return nil
}

// GetCommentIDByVideoID 通过视频ID获取评论ID
func GetCommentIDByVideoID(videoID uint64) ([]uint64, error) {
	key := constant.CommentRedisZSetPrefix + strconv.FormatUint(videoID, 10)
	strList, err := global.CommentRC.ZRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	// ZRange 查不到数据不会返回 redis.Nil
	if len(strList) == 0 {
		return nil, redis.Nil
	}
	result := make([]uint64, 0, len(strList))
	for _, str := range strList {
		num, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}
	return result, nil
}

func SetCommentInfo(ci *CommentInfo) error {
	key := constant.CommentInfoRedisPrefix + strconv.FormatUint(ci.ID, 10)
	ciJSON, err := json.Marshal(ci)
	if err != nil {
		return err
	}
	return global.CommentRC.Set(key, ciJSON,
		constant.CommentInfoExpiration+time.Duration(rand.Intn(200))*time.Minute).Err()
}

func GetCommentInfo(commentID uint64) (*CommentInfo, error) {
	key := constant.CommentInfoRedisPrefix + strconv.FormatUint(commentID, 10)
	result, err := global.CommentRC.Get(key).Result()
	if err != nil {
		return nil, err
	}
	ci := &CommentInfo{}
	err = json.Unmarshal([]byte(result), ci)
	if err != nil {
		return nil, err
	}

	global.CommentRC.Expire(key, constant.CommentInfoExpiration+time.Duration(rand.Intn(200))*time.Minute)

	return ci, nil
}
