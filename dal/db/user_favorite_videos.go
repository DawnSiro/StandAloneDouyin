package db

import (
	"douyin/biz/model/api"
	"douyin/constant"
)

type UserFavoriteVideos struct {
	UserId  uint64 `json:"user_id"`
	VideoId uint64 `json:"video_id"`
}

func (n *UserFavoriteVideos) TableName() string {
	return constant.UserFavoriteVideosTableName
}

func Like(userId uint64, videoId uint64) (uint, error) {
	userFavoriteVideos := &UserFavoriteVideos{
		UserId:  userId,
		VideoId: videoId,
	}

	//is database has this data?
	userFavoriteVideosTemp := &UserFavoriteVideos{}
	result := DB.Where("user_id", userId).Where("video_id", videoId).Find(userFavoriteVideosTemp)
	if result.RowsAffected != 0 {
		return 0, nil
	}

	if err := DB.Create(userFavoriteVideos).Error; err != nil {
		return 0, err
	}
	return 1, nil
}

func UnLike(userId uint64, videoId uint64) (uint, error) {
	userFavoriteVideos := &UserFavoriteVideos{
		UserId:  userId,
		VideoId: videoId,
	}
	//is database has this data?
	userFavoriteVideosTemp := &UserFavoriteVideos{}
	result := DB.Where("user_id", userId).Where("video_id", videoId).Find(userFavoriteVideosTemp)
	if result.RowsAffected != 1 {
		return 0, nil
	}

	if err := DB.Delete(userFavoriteVideos).Error; err != nil {
		return 0, err
	}
	return 1, nil
}

func SelectFavoriteVideoListByUserId(userId uint64) ([]*api.Video, error) {
	resultList := make([]*api.Video, 0)
	userFavoriteVideosList := new([]*UserFavoriteVideos)
	if err := DB.Where("user_id", userId).Find(&userFavoriteVideosList).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(*userFavoriteVideosList); i++ {
		video := &Video{}
		user := &User{}
		DB.Where("id", (*userFavoriteVideosList)[i].VideoId).Find(&video)
		DB.Where("id", video.AuthorID).Find(&user)
		//TODO:Avatar not put
		userTemp := &api.User{
			ID:            int64(user.ID),
			Name:          user.Username,
			FollowCount:   &user.FollowingCount,
			FollowerCount: &user.FollowerCount,
			IsFollow:      IsFollow(userId, uint64(video.AuthorID)),
			Avatar:        "",
		}
		videoTemp := &api.Video{
			ID:            int64(video.ID),
			Author:        userTemp,
			PlayURL:       video.PlayUrl,
			CoverURL:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    true,
			Title:         video.Title,
		}
		resultList = append(resultList, videoTemp)
	}

	return resultList, nil
}
