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
func MSelectFeedVideoDataListByUserID_hotvalue(maxVideoNum int, latestTime *int64, userID uint64) ([]*VideoData, error) {

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


// Define a type for the Item Similarity Matrix.
type ItemSimilarityMatrix map[uint64]map[uint64]float64

// Define a type for User Interactions (e.g., user ratings).
type UserInteractions map[uint64]float64

func GetUserInteractions(maxVideoNum int, latestTime *int64, userID uint64) (UserInteractions, error) {
	// Initialize an empty map to store user interactions (ratings).
	interactions := make(UserInteractions)

	// Query the database to retrieve videos that the user has interacted with.
	videos, err := MSelectFeedVideoDataListByUserID(maxVideoNum, latestTime, userID)
	if err != nil {
		return nil, err
	}

	// Convert video interactions (e.g., favorite_count) into ratings.
	for _, video := range videos {
		// You can customize how you want to calculate ratings based on interactions.
		// Here, we're using favorite_count as a rating.
		rating := float64(video.FavoriteCount)

		// Store the rating in the interactions map with the video ID as the key.
		interactions[video.VID] = rating
	}

	return interactions, nil
}

// CalculateItemSimilarities calculates item-item cosine similarity scores.
func CalculateItemSimilarities(interactions UserInteractions) ItemSimilarityMatrix {
	// Initialize the Item Similarity Matrix.
	similarityMatrix := make(ItemSimilarityMatrix)

	// Iterate through each item pair to calculate their similarity.
	for itemID1 := range interactions {
		similarityMatrix[itemID1] = make(map[uint64]float64)
		for itemID2 := range interactions {
			if itemID1 != itemID2 {
				// Calculate the cosine similarity between item1 and item2.
				similarity := cosineSimilarity(interactions[itemID1], interactions[itemID2])

				// Store the similarity in the matrix.
				similarityMatrix[itemID1][itemID2] = similarity
			}
		}
	}

	return similarityMatrix
}

// FindSimilarItems finds similar items for a given item based on the Item Similarity Matrix.
func FindSimilarItems(itemID uint64, similarityMatrix ItemSimilarityMatrix) map[uint64]float64 {
	// Initialize a map to store similar items and their similarity scores.
	similarItems := make(map[uint64]float64)

	// Iterate through items in the similarity matrix.
	for otherItemID, similarityScore := range similarityMatrix[itemID] {
		// Exclude the same item and items with negative or zero similarity.
		if otherItemID != itemID && similarityScore > 0 {
			similarItems[otherItemID] = similarityScore
		}
	}

	return similarItems
}

// Cosine similarity calculation.
func cosineSimilarity(vector1, vector2 float64) float64 {
	// Calculate the dot product of two vectors.
	dotProduct := vector1 * vector2

	// Calculate the magnitude (Euclidean norm) of each vector.
	magnitude1 := math.Sqrt(vector1 * vector1)
	magnitude2 := math.Sqrt(vector2 * vector2)

	// Calculate the cosine similarity.
	similarity := dotProduct / (magnitude1 * magnitude2)

	return similarity
}

// SortAndSelectTopItems sorts and selects the top N items based on scores.
func SortAndSelectTopItems(scores map[uint64]float64, maxItems int) []uint64 {
	// Create a slice of items sorted by their scores.
	sortedItems := make([]uint64, 0, len(scores))
	for itemID := range scores {
		sortedItems = append(sortedItems, itemID)
	}

	sort.Slice(sortedItems, func(i, j int) bool {
		// Sort items in descending order of scores.
		return scores[sortedItems[i]] > scores[sortedItems[j]]
	})

	// Select up to maxItems items or all items if there are fewer than maxItems.
	numItems := len(sortedItems)
	if numItems > maxItems {
		return sortedItems[:maxItems]
	}

	return sortedItems
}

func FetchVideoDataForItems(itemIDs []uint64, numToReturn int, latestTime *int64, userID uint64, topItemIDs []uint64) ([]*VideoData, error) {
	// Create a slice to store the fetched video data.
	videoDataSlice := make([]*VideoData, 0)

	// Implement your database query to fetch video data based on the provided condition.
	// Example database query:
	query := `
        SELECT
            v.id AS vid,
            v.play_url,
            v.cover_url,
            v.favorite_count,
            v.comment_count,
            v.title,
            u.id AS uid,
            u.username,
            u.following_count,
            u.follower_count,
            u.avatar,
            u.background_image,
            u.signature,
            u.total_favorited,
            u.work_count,
            u.favorite_count AS user_favorite_count,
            IF(r.is_deleted = ?, TRUE, FALSE) AS is_follow,
            IF(ufv.is_deleted = ?, TRUE, FALSE) AS is_favorite
        FROM
            video AS v
        RIGHT JOIN
            user AS u ON u.id = v.author_id
        LEFT JOIN
            relation AS r ON r.to_user_id = u.id AND r.user_id = ?
        LEFT JOIN
            user_favorite_video AS ufv ON v.id = ufv.video_id AND ufv.user_id = ?
        WHERE
            v.id IN (?)
            AND v.publish_time < ?
        LIMIT ?
    `

	// Execute the query and scan the results.
	rows, err := global.DB.Raw(query,
		constant.DataNotDeleted,
		constant.DataNotDeleted,
		userID, // Replace with your user ID
		userID, // Replace with your user ID
		itemIDs,
		time.UnixMilli(*latestTime),
		numToReturn,
	).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the query result and populate the videoDataSlice.
	for rows.Next() {
		var video VideoData
		if err := rows.Scan(
			&video.VID,
			&video.PlayURL,
			&video.CoverURL,
			&video.FavoriteCount,
			&video.CommentCount,
			&video.Title,
			&video.UID,
			&video.Username,
			&video.FollowCount,
			&video.FollowerCount,
			&video.Avatar,
			&video.BackgroundImage,
			&video.Signature,
			&video.TotalFavorited,
			&video.WorkCount,
			&video.UserFavoriteCount,
			&video.IsFollow,
			&video.IsFavorite,
		); err != nil {
			return nil, err
		}
		videoDataSlice = append(videoDataSlice, &video)
	}

	// Filter the results based on topItemIDs.
	filteredVideoDataSlice := make([]*VideoData, 0)
	for _, video := range videoDataSlice {
		for _, topItemID := range topItemIDs {
			if video.VID == topItemID {
				filteredVideoDataSlice = append(filteredVideoDataSlice, video)
				break
			}
		}
	}

	// Return only the first numToReturn elements from filteredVideoDataSlice.
	if numToReturn >= len(filteredVideoDataSlice) {
		return filteredVideoDataSlice, nil
	}
	return filteredVideoDataSlice[:numToReturn], nil
}

// GenerateItemCFRecommendations generates recommendations using Item CF.
func GenerateItemCFRecommendations(maxVideoNum int, latestTime *int64, userID uint64) ([]*VideoData, error) {
	// Get user interactions (ratings) from the database.
	interactions, err := GetUserInteractions(maxVideoNum, latestTime, userID)
	if err != nil {
		return nil, err
	}

	// Calculate item-item similarity scores.
	similarityMatrix := CalculateItemSimilarities(interactions)

	// Initialize a map to store item recommendations and their scores.
	recommendations := make(map[uint64]float64)

	// Iterate through items the user has interacted with.
	for itemID, userRating := range interactions {
		// Find similar items to the current item based on item similarity scores.
		similarItems := FindSimilarItems(itemID, similarityMatrix)

		// Iterate through similar items and update recommendations.
		for similarItemID, similarityScore := range similarItems {
			// Skip items the user has already interacted with.
			if _, exists := interactions[similarItemID]; exists {
				continue
			}

			// Update or add to the recommendations based on similarity and user rating.
			recommendations[similarItemID] += similarityScore * userRating
		}
	}

	// Sort and select the top recommendations.
	topItemIDs := SortAndSelectTopItems(recommendations, maxVideoNum)

	// Fetch video data for the recommended items.
	recommendedVideoData, err := FetchVideoDataForItems(topItemIDs, maxVideoNum, latestTime, userID, topItemIDs)
	if err != nil {
		return nil, err
	}

	return recommendedVideoData, nil
}
