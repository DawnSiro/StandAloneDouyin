package rdb

import (
	"douyin/pkg/constant"
	"douyin/pkg/global"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"strconv"
	"time"
)

type Video struct {
	ID uint64 `json:"id"`
}

// VideoInfo 视频固定的信息
type VideoInfo struct {
	PublishTime time.Time
	// 有单独获取的需求，故 AuthorID 单独使用 Redis String 进行缓存
	AuthorID uint64 `gorm:"not null" json:"author_id"`
	PlayURL  string
	CoverURL string
	Title    string
}

type VideoCount struct {
	FavoriteCount int64 `gorm:"default:0;not null" json:"favorite_count"`
	CommentCount  int64 `gorm:"default:0;not null" json:"comment_count"`
}

// LoadFavoriteVideoID 加载用户点赞视频列表到 Redis
func LoadFavoriteVideoID(userID uint64, videoIDs []uint64) error {
	key := constant.FavoriteVideoIDRedisSetPrefix + strconv.FormatUint(userID, 10)
	err := global.VideoRC.ZAdd(key).Err()
	if err != nil {
		hlog.Error("rdb.video.FavoriteVideo err:", err.Error())
		return err
	}
	return nil
}

// AddFavoriteVideoID 添加 Redis 用户点赞视频 ID
func AddFavoriteVideoID(userID, videoID uint64) error {
	key := constant.FavoriteVideoIDRedisSetPrefix + strconv.FormatUint(userID, 10)
	err := global.VideoRC.SAdd(key, videoID).Err()
	if err != nil {
		hlog.Error("rdb.video.FavoriteVideo err:", err.Error())
		return err
	}
	return nil
}

// DelFavoriteVideoID 删除 Redis 用户点赞视频 ID
func DelFavoriteVideoID(userID, videoID uint64) error {
	key := constant.FavoriteVideoIDRedisSetPrefix + strconv.FormatUint(userID, 10)
	err := global.VideoRC.SRem(key, videoID).Err()
	if err != nil {
		hlog.Error("rdb.video.DelFavoriteVideo err:", err.Error())
		return err
	}
	return nil
}

// SetAuthorID 通过视频ID获取
func SetAuthorID(authorID, videoID uint64) error {
	key := constant.VideoAuthorIDRedisPrefix + strconv.FormatUint(videoID, 10)
	return global.VideoRC.Set(key, authorID, constant.VideoAuthorIDExpiration).Err()
}

// GetAuthorID 通过视频ID获取作者ID
func GetAuthorID(videoID uint64) (uint64, error) {
	key := constant.VideoAuthorIDRedisPrefix + strconv.FormatUint(videoID, 10)
	authorIDString, err := global.VideoRC.Get(key).Result()
	if err != nil {
		return 0, nil
	}
	return strconv.ParseUint(authorIDString, 10, 64)
}
