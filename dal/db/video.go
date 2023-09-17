package db

import (
	"time"

	"douyin/dal/model"
	"douyin/pkg/constant"
	"douyin/pkg/global"

	"gorm.io/gorm"
)

// CreateVideo 新增视频，返回的 videoID 是为了将 videoID 放入布隆过滤器
func CreateVideo(video *model.Video) (uint64, error) {
	var videoID uint64
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		//从这里开始，应该使用 tx 而不是 db（tx 是 Transaction 的简写）
		u := &model.User{ID: video.AuthorID}
		err := tx.Select("work_count").First(u).Error
		if err != nil {
			return err
		}
		err = tx.Model(u).Update("work_count", u.WorkCount+1).Error
		if err != nil {
			return err
		}
		err = tx.Create(video).Error
		if err != nil {
			return err
		}
		videoID = video.ID
		// 返回 nil 提交事务
		return nil
	})
	if err != nil {
		return 0, err
	}
	return videoID, nil
}

// MGetVideos multiple get list of videos info
func MGetVideos(maxVideoNum int, latestTime *int64) ([]*model.Video, error) {
	res := make([]*model.Video, 0)

	if latestTime == nil || *latestTime == 0 {
		// 这里和文档里说得不一样，实际客户端传的是毫秒
		currentTime := time.Now().UnixMilli()
		latestTime = &currentTime
	}

	// TODO 设计简单的推荐算法，比如关注的 UP 发了视频，会优先推送
	if err := global.DB.Where("publish_time < ?", time.UnixMilli(*latestTime)).Limit(maxVideoNum).
		Order("publish_time desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func GetVideosByAuthorID(userID uint64) ([]*model.Video, error) {
	res := make([]*model.Video, 0)
	err := global.DB.Find(&res, "author_id = ? ", userID).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SelectAuthorIDByVideoID(videoID uint64) (uint64, error) {
	video := &model.Video{
		ID: videoID,
	}

	err := global.DB.First(&video).Error
	if err != nil {
		return 0, err
	}
	return video.AuthorID, nil
}

func UpdateVideoFavoriteCount(videoID uint64, favoriteCount uint64) (int64, error) {
	video := &model.Video{
		ID: videoID,
	}

	if err := global.DB.Model(&video).Update("favorite_count", favoriteCount).Error; err != nil {
		return 0, err
	}
	return video.FavoriteCount, nil
}

// IncreaseVideoFavoriteCount increase 1
func IncreaseVideoFavoriteCount(videoID uint64) (int64, error) {
	video := &model.Video{
		ID: videoID,
	}
	err := global.DB.First(&video).Error
	if err != nil {
		return 0, err
	}
	if err := global.DB.Model(&video).Update("favorite_count", video.FavoriteCount+1).Error; err != nil {
		return 0, err
	}
	return video.CommentCount, nil
}

// DecreaseVideoFavoriteCount decrease 1
func DecreaseVideoFavoriteCount(videoID uint64) (int64, error) {
	video := &model.Video{
		ID: videoID,
	}
	err := global.DB.First(&video).Error
	if err != nil {
		return 0, err
	}
	if err := global.DB.Model(&video).Update("favorite_count", video.FavoriteCount-1).Error; err != nil {
		return 0, err
	}
	return video.CommentCount, nil
}

func UpdateCommentCount(videoID uint64, commentCount uint64) (int64, error) {
	video := &model.Video{
		ID: videoID,
	}

	if err := global.DB.Model(&video).Update("comment_count", commentCount).Error; err != nil {
		return 0, err
	}
	return video.FavoriteCount, nil
}

// IncreaseCommentCount increase 1
func IncreaseCommentCount(videoID uint64) (int64, error) {
	video := &model.Video{
		ID: videoID,
	}

	err := global.DB.First(&video).Error
	if err != nil {
		return 0, err
	}

	if err := global.DB.Model(&video).Update("comment_count", video.CommentCount+1).Error; err != nil {
		return 0, err
	}
	return video.CommentCount, nil
}

// DecreaseCommentCount decrease  1
func DecreaseCommentCount(videoID uint64) (int64, error) {
	video := &model.Video{
		ID: videoID,
	}
	err := global.DB.First(&video).Error
	if err != nil {
		return 0, err
	}
	if err := global.DB.Model(&video).Update("comment_count", video.CommentCount-1).Error; err != nil {
		return 0, err
	}
	return video.CommentCount, nil
}

func SelectVideoFavoriteCountByVideoID(videoID uint64) (int64, error) {
	video := &model.Video{
		ID: videoID,
	}

	err := global.DB.Select("favorite_count").First(&video).Error
	if err != nil {
		return 0, err
	}
	return video.FavoriteCount, nil
}

func SelectCommentCountByVideoID(videoID uint64) (int64, error) {
	video := &model.Video{
		ID: videoID,
	}

	err := global.DB.Select("comment_count").First(&video).Error
	if err != nil {
		return 0, err
	}
	return video.CommentCount, nil
}

// SelectFavoriteVideoDataListByUserID 查询点赞视频列表
func SelectFavoriteVideoDataListByUserID(userID, selectUserID uint64) ([]*model.VideoData, error) {
	res := make([]*model.VideoData, 0)
	sqlQuery := global.DB.Select("ufv.video_id").
		Table("user_favorite_video AS ufv").
		Where("ufv.user_id = ? AND ufv.is_deleted = ?", selectUserID, constant.DataNotDeleted)
	err := global.DB.Select("v.id AS vid, v.play_url, v.cover_url, v.favorite_count, v.comment_count, v.title,"+
		"u.id AS uid, u.username, u.following_count, u.follower_count, u.avatar,"+
		"u.background_image, u.signature, u.total_favorited, u.work_count, u.favorite_count as user_favorite_count,"+
		"IF(r.is_deleted = ?, TRUE, FALSE) AS is_follow, IF(ufv.is_deleted = ?, TRUE, FALSE) AS is_favorite",
		constant.DataNotDeleted, constant.DataNotDeleted).Table("user AS u").
		Joins("RIGHT JOIN video AS v ON u.id = v.author_id").
		Joins("LEFT JOIN relation AS r ON r.to_user_id = u.id AND r.user_id = ?", userID).
		Joins("LEFT JOIN user_favorite_video AS ufv ON v.id = ufv.video_id AND ufv.user_id = ?", userID).
		Where("v.id IN (?)", sqlQuery).Order("v.publish_time DESC").Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

// MSelectFeedVideoDataListByUserID 查询Feed视频列表
func MSelectFeedVideoDataListByUserID(maxVideoNum int, latestTime *int64, userID uint64) ([]*model.VideoData, error) {
	res := make([]*model.VideoData, 0)
	if latestTime == nil || *latestTime == 0 {
		// 这里和文档里说得不一样，实际客户端传的是毫秒
		currentTime := time.Now().UnixMilli()
		latestTime = &currentTime
	}
	err := global.DB.Select("v.id AS vid, v.play_url, v.cover_url, v.favorite_count, v.comment_count, v.title,"+
		"u.id AS uid, u.username, u.following_count, u.follower_count, u.avatar,"+
		"u.background_image, u.signature, u.total_favorited, u.work_count, u.favorite_count as user_favorite_count,"+
		"IF(r.is_deleted = ?, TRUE, FALSE) AS is_follow, IF(ufv.is_deleted = ?, TRUE, FALSE) AS is_favorite",
		constant.DataNotDeleted, constant.DataNotDeleted).Table("user AS u").
		Joins("RIGHT JOIN video AS v ON u.id = v.author_id").
		Joins("LEFT JOIN relation AS r ON r.to_user_id = u.id AND r.user_id = 6", userID).
		Joins("LEFT JOIN user_favorite_video AS ufv ON v.id = ufv.video_id AND ufv.user_id = ?", userID).
		Where("v.publish_time < ?", time.UnixMilli(*latestTime)).Limit(maxVideoNum).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SelectPublishTimeByVideoID(videoID uint64) (int64, error) {
	var publishTime time.Time
	err := global.DB.Select("publish_time").Model(&model.Video{ID: videoID}).First(&publishTime).Error
	if err != nil {
		return 0, err
	}
	return publishTime.UnixMilli(), nil
}

// SelectPublishVideoDataListByUserID 发布视频列表
func SelectPublishVideoDataListByUserID(userID, selectUserID uint64) ([]*model.VideoData, error) {
	res := make([]*model.VideoData, 0)
	err := global.DB.Select("v.id AS vid, v.play_url, v.cover_url, v.favorite_count, v.comment_count, v.title,"+
		"u.id AS uid, u.username, u.following_count, u.follower_count, u.avatar,"+
		"u.background_image, u.signature, u.total_favorited, u.work_count, u.favorite_count as user_favorite_count,"+
		"IF(r.is_deleted = ?, TRUE, FALSE) AS is_follow, IF(ufv.is_deleted = ?, TRUE, FALSE) AS is_favorite",
		constant.DataNotDeleted, constant.DataNotDeleted).Table("user AS u").
		Joins("RIGHT JOIN video AS v ON u.id = v.author_id").
		Joins("LEFT JOIN relation AS r ON r.to_user_id = u.id AND r.user_id = ?", userID).
		Joins("LEFT JOIN user_favorite_video AS ufv ON v.id = ufv.video_id AND ufv.user_id = ?", userID).
		Where("v.author_id = ?", selectUserID).Order("v.publish_time DESC").Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
