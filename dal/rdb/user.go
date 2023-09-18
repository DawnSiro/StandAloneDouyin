package rdb

import (
	"douyin/dal/model"
	"douyin/pkg/constant"
	"douyin/pkg/global"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"math/rand"
	"strconv"
	"time"
)

// UserInfo 用户信息中基本信息，不做更改或更改频率较低
type UserInfo struct {
	ID              uint64 // 自增主键
	Username        string // 用户名
	Password        string // 用户密码
	Avatar          string // 用户头像
	BackgroundImage string // 用户个人页顶部大图
	Signature       string // 个人简介
}

// UserInfoCount 用户信息中的计数信息，需要频繁更新
type UserInfoCount struct {
	FollowingCount      int64 //关注数
	FollowerCount       int64 //粉丝数
	TotalFavoritedCount int64 //获赞数量
	WorkCount           int64 //作品数量
	FavoriteCount       int64 //点赞数量
}

// SetUserInfo 设置用户信息
func SetUserInfo(user *model.User) error {
	// 拆分成两个数据结构，不做更改或更改频率较低放在 UserInfo 中
	userInfo := &UserInfo{
		ID:              user.ID,
		Username:        user.Username,
		Password:        user.Password,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
	}

	// 进行序列化
	userInfoJSON, err := json.Marshal(userInfo)
	if err != nil {
		hlog.Error("dal.rdb.user.SetUserInfo err:", err.Error())
		return err
	}

	// 开启管道
	pipeline := global.UserRC.Pipeline()

	// 设置 UserInfo 的 JSON 缓存
	infoKey := constant.UserInfoRedisPrefix + strconv.FormatUint(user.ID, 10)
	err = pipeline.Set(infoKey, userInfoJSON,
		constant.UserInfoExpiration+time.Duration(rand.Intn(200))*time.Minute).Err()
	if err != nil {
		hlog.Error("dal.rdb.user.SetUserInfo err:", err.Error())
		return err
	}

	// 需要计数更新的部分放入 UserInfoCount 的缓存中
	userInfoCount := &UserInfoCount{
		FollowingCount:      user.FollowingCount,
		FollowerCount:       user.FollowerCount,
		TotalFavoritedCount: user.TotalFavorited,
		WorkCount:           user.WorkCount,
		FavoriteCount:       user.FavoriteCount,
	}

	infoCountKey := constant.UserInfoCountHashRedisPrefix + strconv.FormatUint(user.ID, 10)
	// 使用 MSet 进行批量设置
	err = pipeline.HMSet(infoCountKey, map[string]interface{}{
		constant.FollowCountRedisFiled:    userInfoCount.FollowingCount,
		constant.FanCountRedisFiled:       userInfoCount.FollowerCount,
		constant.TotalFavoritedRedisFiled: userInfoCount.TotalFavoritedCount,
		constant.WorkCountRedisFiled:      userInfoCount.WorkCount,
		constant.FavoriteCountRedisFiled:  userInfoCount.FavoriteCount,
	}).Err()
	if err != nil {
		hlog.Error("dal.rdb.user.SetUserInfo err:", err.Error())
		return err
	}
	err = pipeline.Expire(infoCountKey,
		constant.UserInfoExpiration+time.Duration(rand.Intn(200))*time.Minute).Err()
	if err != nil {
		hlog.Error("dal.rdb.user.SetUserInfo err:", err.Error())
		return err
	}
	// 执行管道中的命令
	_, err = pipeline.Exec()
	if err != nil {
		hlog.Error("dal.rdb.user.SetUserInfo err:", err.Error())
		return err
	}
	return nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(userID uint64) (*model.User, error) {
	// 获取用户信息中基本信息，不做更改或更改频率较低的 UserInfo
	infoKey := constant.UserInfoRedisPrefix + strconv.FormatUint(userID, 10)
	userInfoJSON, err := global.UserRC.Get(infoKey).Result()
	if err != nil {
		hlog.Error("dal.rdb.user.GetUserInfo err:", err.Error())
		return nil, err
	}
	userinfo := &UserInfo{}
	err = json.Unmarshal([]byte(userInfoJSON), userinfo)
	if err != nil {
		hlog.Error("dal.rdb.user.GetUserInfo err:", err.Error())
		return nil, err
	}

	// 获取用户信息中的计数信息，需要频繁更新的 UserInfoCount
	infoCountKey := constant.UserInfoCountHashRedisPrefix + strconv.FormatUint(userID, 10)
	countMap, err := global.UserRC.HGetAll(infoCountKey).Result()
	if err != nil {
		hlog.Error("dal.rdb.user.GetUserInfo err:", err.Error())
		return nil, err
	}
	// 进行参数解析
	followingCount, err := strconv.ParseInt(countMap[constant.FollowCountRedisFiled], 10, 64)
	followerCount, err := strconv.ParseInt(countMap[constant.FanCountRedisFiled], 10, 64)
	totalFavoritedCount, err := strconv.ParseInt(countMap[constant.TotalFavoritedRedisFiled], 10, 64)
	workCount, err := strconv.ParseInt(countMap[constant.WorkCountRedisFiled], 10, 64)
	favoriteCount, err := strconv.ParseInt(countMap[constant.FollowCountRedisFiled], 10, 64)
	if err != nil {
		hlog.Error("dal.rdb.user.GetUserInfo err:", err.Error())
		return nil, err
	}
	userInfoCount := &UserInfoCount{
		FollowingCount:      followingCount,
		FollowerCount:       followerCount,
		TotalFavoritedCount: totalFavoritedCount,
		WorkCount:           workCount,
		FavoriteCount:       favoriteCount,
	}

	return &model.User{
		ID:              userinfo.ID,
		Username:        userinfo.Username,
		FollowingCount:  userInfoCount.FollowingCount,
		FollowerCount:   userInfoCount.FollowerCount,
		Avatar:          userinfo.Avatar,
		BackgroundImage: userinfo.BackgroundImage,
		Signature:       userinfo.Signature,
		TotalFavorited:  userInfoCount.TotalFavoritedCount,
		WorkCount:       userInfoCount.WorkCount,
		FavoriteCount:   userInfoCount.FavoriteCount,
	}, nil
}

// ExpireUserInfo 更新用户信息过期时间
func ExpireUserInfo(userID uint64) error {
	infoKey := constant.UserInfoRedisPrefix + strconv.FormatUint(userID, 10)
	infoCountKey := constant.UserInfoRedisPrefix + strconv.FormatUint(userID, 10)
	// 使用管道加速
	pipeline := global.UserRC.Pipeline()
	// 更新过期时间
	err := pipeline.Expire(infoKey, constant.UserInfoExpiration).Err()
	if err != nil {
		hlog.Error("dal.rdb.user.ExpireUserInfo err:", err.Error())
		return err
	}
	err = pipeline.Expire(infoCountKey, constant.UserInfoExpiration).Err()
	if err != nil {
		hlog.Error("dal.rdb.user.ExpireUserInfo err:", err.Error())
		return err
	}
	return nil
}
