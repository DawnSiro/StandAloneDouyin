package db

import (
	"douyin/pkg/constant"
	"douyin/pkg/global"
	"math"
	"sort"
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID            uint64    `json:"id"`
	PublishTime   time.Time `gorm:"not null" json:"publish_time"`
	AuthorID      uint64    `gorm:"not null" json:"author_id"`
	PlayURL       string    `gorm:"type:varchar(255);not null" json:"play_url"`
	CoverURL      string    `gorm:"type:varchar(255);not null" json:"cover_url"`
	FavoriteCount int64     `gorm:"default:0;not null" json:"favorite_count"`
	CommentCount  int64     `gorm:"default:0;not null" json:"comment_count"`
	Title         string    `gorm:"type:varchar(63);not null" json:"title"`
}

func (n *Video) TableName() string {
	return constant.VideoTableName
}

func CreateVideo(video *Video) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		//从这里开始，应该使用 tx 而不是 db（tx 是 Transaction 的简写）
		u := &User{ID: video.AuthorID}
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
		// 返回 nil 提交事务
		return nil
	})
}

// MGetVideos multiple get list of videos info
func MGetVideos(maxVideoNum int, latestTime *int64) ([]*Video, error) {
	res := make([]*Video, 0)

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

func GetVideosByAuthorID(userID uint64) ([]*Video, error) {
	res := make([]*Video, 0)
	err := global.DB.Find(&res, "author_id = ? ", userID).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SelectAuthorIDByVideoID(videoID uint64) (uint64, error) {
	video := &Video{
		ID: videoID,
	}

	err := global.DB.First(&video).Error
	if err != nil {
		return 0, err
	}
	return video.AuthorID, nil
}

func UpdateVideoFavoriteCount(videoID uint64, favoriteCount uint64) (int64, error) {
	video := &Video{
		ID: videoID,
	}

	if err := global.DB.Model(&video).Update("favorite_count", favoriteCount).Error; err != nil {
		return 0, err
	}
	return video.FavoriteCount, nil
}

// IncreaseVideoFavoriteCount increase 1
func IncreaseVideoFavoriteCount(videoID uint64) (int64, error) {
	video := &Video{
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
	video := &Video{
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
	video := &Video{
		ID: videoID,
	}

	if err := global.DB.Model(&video).Update("comment_count", commentCount).Error; err != nil {
		return 0, err
	}
	return video.FavoriteCount, nil
}

// IncreaseCommentCount increase 1
func IncreaseCommentCount(videoID uint64) (int64, error) {
	video := &Video{
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
	video := &Video{
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
	video := &Video{
		ID: videoID,
	}

	err := global.DB.Select("favorite_count").First(&video).Error
	if err != nil {
		return 0, err
	}
	return video.FavoriteCount, nil
}

func SelectCommentCountByVideoID(videoID uint64) (int64, error) {
	video := &Video{
		ID: videoID,
	}

	err := global.DB.Select("comment_count").First(&video).Error
	if err != nil {
		return 0, err
	}
	return video.CommentCount, nil
}

type VideoData struct {
	// 没有ID的话，貌似第一个字段会被识别为主键ID
	VID               uint64 `gorm:"column:vid"`
	PlayURL           string
	CoverURL          string
	FavoriteCount     int64
	CommentCount      int64
	IsFavorite        bool
	Title             string
	UID               uint64
	Username          string
	FollowCount       int64
	FollowerCount     int64
	IsFollow          bool
	Avatar            string
	BackgroundImage   string
	Signature         string
	TotalFavorited    int64
	WorkCount         int64
	UserFavoriteCount int64
	HotValue          float64
	PublishTime       time.Time
}

// SelectFavoriteVideoDataListByUserID 查询点赞视频列表
func SelectFavoriteVideoDataListByUserID(userID, selectUserID uint64) ([]*VideoData, error) {
	res := make([]*VideoData, 0)
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
func MSelectFeedVideoDataListByUserID(maxVideoNum int, latestTime *int64, userID uint64) ([]*VideoData, error) {
	res := make([]*VideoData, 0)
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
	err := global.DB.Select("publish_time").Model(&Video{ID: videoID}).First(&publishTime).Error
	if err != nil {
		return 0, err
	}
	return publishTime.UnixMilli(), nil
}

// SelectPublishVideoDataListByUserID 发布视频列表
func SelectPublishVideoDataListByUserID(userID, selectUserID uint64) ([]*VideoData, error) {
	res := make([]*VideoData, 0)
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

func MSelectFeedVideoDataListByUserID2(maxVideoNum int, latestTime *int64, userID uint64) ([]*VideoData, error) {
	res := make([]*VideoData, 0)
	if latestTime == nil || *latestTime == 0 {
		// This part remains the same as in your original function.
		currentTime := time.Now().UnixMilli()
		latestTime = &currentTime
	}

	// Fetch video data from the database without any filtering based on hotness.
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

	// Calculate and set the hotness for each video.
	decayFactor := 0.95
	for _, video := range res {
		hotValue := calculateHotValue(video, *latestTime, decayFactor)
		video.HotValue = hotValue
	}

	// Apply hotness-based filtering and return the filtered videos.
	recommendedVideoData := filterByHotness(res, maxVideoNum)

	return recommendedVideoData, nil
}

func calculateHotValue(video *VideoData, currentTime int64, decayFactor float64) float64 {
	// Fetch the publish time for the video
	publishTime, err := SelectPublishTimeByVideoID(video.VID)
	if err != nil {
		// Handle error (you can log it or return a default value)
		return 0.0
	}

	// Calculate the elapsed time in hours
	elapsedHours := float64(currentTime-publishTime) / float64(time.Hour)

	// Calculate the original hotness value based on likes and comments (adjust weights as needed)
	originalHotValue := float64(video.FavoriteCount)*0.4 + float64(video.CommentCount)*0.5

	// Calculate the decayed hotness value based on the elapsed time and decay factor
	decayedHotValue := originalHotValue * math.Pow(decayFactor, elapsedHours)

	return decayedHotValue
}

func filterByHotness(videos []*VideoData, maxNum int) []*VideoData {
	// Sort the videos by hotness in descending order.
	sort.SliceStable(videos, func(i, j int) bool {
		return videos[i].HotValue > videos[j].HotValue
	})

	// Determine the number of videos to return (up to the maximum specified).
	numToReturn := len(videos)
	if numToReturn > maxNum {
		numToReturn = maxNum
	}

	// Return the top numToReturn videos based on hotness.
	return videos[:numToReturn]
}
