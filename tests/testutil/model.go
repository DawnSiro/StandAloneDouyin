package testutil

import (
	"encoding/json"
	"io"
	"net/http"
)

func GetDouyinResponse[T any](resp *http.Response) (respData T, err error) {
	var data []byte
	if data, err = io.ReadAll(resp.Body); err == nil {
		err = json.Unmarshal(data, &respData)
	}
	return
}

type UserInfo struct {
	ID              int64  `thrift:"id,1,required" form:"id,required" json:"id,required" query:"id,required"`
	Name            string `thrift:"name,2,required" form:"name,required" json:"name,required" query:"name,required"`
	FollowCount     int64  `thrift:"follow_count,3,required" form:"follow_count,required" json:"follow_count,required" query:"follow_count,required"`
	FollowerCount   int64  `thrift:"follower_count,4,required" form:"follower_count,required" json:"follower_count,required" query:"follower_count,required"`
	IsFollow        bool   `thrift:"is_follow,5,required" form:"is_follow,required" json:"is_follow,required" query:"is_follow,required"`
	Avatar          string `thrift:"avatar,6,required" form:"avatar,required" json:"avatar,required" query:"avatar,required"`
	BackgroundImage string `thrift:"background_image,7,required" form:"background_image,required" json:"background_image,required" query:"background_image,required"`
	Signature       string `thrift:"signature,8,required" form:"signature,required" json:"signature,required" query:"signature,required"`
	TotalFavorited  int64  `thrift:"total_favorited,9,required" form:"total_favorited,required" json:"total_favorited,required" query:"total_favorited,required"`
	WorkCount       int64  `thrift:"work_count,10,required" form:"work_count,required" json:"work_count,required" query:"work_count,required"`
	FavoriteCount   int64  `thrift:"favorite_count,11,required" form:"favorite_count,required" json:"favorite_count,required" query:"favorite_count,required"`
}

type Video struct {
	ID            int64     `thrift:"id,1,required" form:"id,required" json:"id,required" query:"id,required"`
	Author        *UserInfo `thrift:"author,2,required" form:"author,required" json:"author,required" query:"author,required"`
	PlayURL       string    `thrift:"play_url,3,required" form:"play_url,required" json:"play_url,required" query:"play_url,required"`
	CoverURL      string    `thrift:"cover_url,4,required" form:"cover_url,required" json:"cover_url,required" query:"cover_url,required"`
	FavoriteCount int64     `thrift:"favorite_count,5,required" form:"favorite_count,required" json:"favorite_count,required" query:"favorite_count,required"`
	CommentCount  int64     `thrift:"comment_count,6,required" form:"comment_count,required" json:"comment_count,required" query:"comment_count,required"`
	IsFavorite    bool      `thrift:"is_favorite,7,required" form:"is_favorite,required" json:"is_favorite,required" query:"is_favorite,required"`
	Title         string    `thrift:"title,8,required" form:"title,required" json:"title,required" query:"title,required"`
}

type DouyinFeedResponse struct {
	StatusCode int64    `thrift:"status_code,1,required" form:"status_code,required" json:"status_code,required" query:"status_code,required"`
	StatusMsg  *string  `thrift:"status_msg,2,optional" form:"status_msg" json:"status_msg,omitempty" query:"status_msg"`
	VideoList  []*Video `thrift:"video_list,3" form:"video_list" json:"video_list" query:"video_list"`
	NextTime   *int64   `thrift:"next_time,4,optional" form:"next_time" json:"next_time,omitempty" query:"next_time"`
}

type DouyinUserRegisterResponse struct {
	StatusCode int64   `thrift:"status_code,1,required" form:"status_code,required" json:"status_code,required" query:"status_code,required"`
	StatusMsg  *string `thrift:"status_msg,2,optional" form:"status_msg" json:"status_msg,omitempty" query:"status_msg"`
	UserID     int64   `thrift:"user_id,3,required" form:"user_id,required" json:"user_id,required" query:"user_id,required"`
	Token      string  `thrift:"token,4,required" form:"token,required" json:"token,required" query:"token,required"`
}

type DouyinUserLoginResponse struct {
	StatusCode int64   `thrift:"status_code,1,required" form:"status_code,required" json:"status_code,required" query:"status_code,required"`
	StatusMsg  *string `thrift:"status_msg,2,optional" form:"status_msg" json:"status_msg,omitempty" query:"status_msg"`
	UserID     int64   `thrift:"user_id,3,required" form:"user_id,required" json:"user_id,required" query:"user_id,required"`
	Token      string  `thrift:"token,4,required" form:"token,required" json:"token,required" query:"token,required"`
}

type DouyinUserResponse struct {
	StatusCode int64     `thrift:"status_code,1,required" form:"status_code,required" json:"status_code,required" query:"status_code,required"`
	StatusMsg  *string   `thrift:"status_msg,2,optional" form:"status_msg" json:"status_msg,omitempty" query:"status_msg"`
	User       *UserInfo `thrift:"user,3,required" form:"user,required" json:"user,required" query:"user,required"`
}

// type DouyinPublishListResponse struct {
// 	StatusCode int64    `thrift:"status_code,1,required" form:"status_code,required" json:"status_code,required" query:"status_code,required"`
// 	StatusMsg  *string  `thrift:"status_msg,2,optional" form:"status_msg" json:"status_msg,omitempty" query:"status_msg"`
// 	VideoList  []*Video `thrift:"video_list,3" form:"video_list" json:"video_list" query:"video_list"`
// }

type DouyinSimpleResponse struct {
	StatusCode int64   `thrift:"status_code,1,required" form:"status_code,required" json:"status_code,required" query:"status_code,required"`
	StatusMsg  *string `thrift:"status_msg,2,optional" form:"status_msg" json:"status_msg,omitempty" query:"status_msg"`
}

type User struct {
	ID            int64  `thrift:"id,1,required" form:"id,required" json:"id,required" query:"id,required"`
	Name          string `thrift:"name,2,required" form:"name,required" json:"name,required" query:"name,required"`
	FollowCount   *int64 `thrift:"follow_count,3,optional" form:"follow_count" json:"follow_count,omitempty" query:"follow_count"`
	FollowerCount *int64 `thrift:"follower_count,4,optional" form:"follower_count" json:"follower_count,omitempty" query:"follower_count"`
	IsFollow      bool   `thrift:"is_follow,5,required" form:"is_follow,required" json:"is_follow,required" query:"is_follow,required"`
	Avatar        string `thrift:"avatar,6,required" form:"avatar,required" json:"avatar,required" query:"avatar,required"`
}

type Comment struct {
	ID         int64  `thrift:"id,1,required" form:"id,required" json:"id,required" query:"id,required"`
	User       *User  `thrift:"user,2,required" form:"user,required" json:"user,required" query:"user,required"`
	Content    string `thrift:"content,3,required" form:"content,required" json:"content,required" query:"content,required"`
	CreateDate string `thrift:"create_date,4,required" form:"create_date,required" json:"create_date,required" query:"create_date,required"`
}

type DouyinCommentActionResponse struct {
	StatusCode int64    `thrift:"status_code,1,required" form:"status_code,required" json:"status_code,required" query:"status_code,required"`
	StatusMsg  *string  `thrift:"status_msg,2,optional" form:"status_msg" json:"status_msg,omitempty" query:"status_msg"`
	Comment    *Comment `thrift:"comment,3,optional" form:"comment" json:"comment,omitempty" query:"comment"`
}

type DouyinCommentListResponse struct {
	StatusCode  int64      `thrift:"status_code,1,required" form:"status_code,required" json:"status_code,required" query:"status_code,required"`
	StatusMsg   *string    `thrift:"status_msg,2,optional" form:"status_msg" json:"status_msg,omitempty" query:"status_msg"`
	CommentList []*Comment `thrift:"comment_list,3" form:"comment_list" json:"comment_list" query:"comment_list"`
}

type DouyinFavoriteListResponse struct {
	StatusCode int64    `thrift:"status_code,1,required" form:"status_code,required" json:"status_code,required" query:"status_code,required"`
	StatusMsg  *string  `thrift:"status_msg,2,optional" form:"status_msg" json:"status_msg,omitempty" query:"status_msg"`
	VideoList  []*Video `thrift:"video_list,3" form:"video_list" json:"video_list" query:"video_list"`
}

type DouyinRelationFollowListResponse struct {
	StatusCode int64   `thrift:"status_code,1,required" form:"status_code,required" json:"status_code,required" query:"status_code,required"`
	StatusMsg  *string `thrift:"status_msg,2,optional" form:"status_msg" json:"status_msg,omitempty" query:"status_msg"`
	UserList   []*User `thrift:"user_list,3" form:"user_list" json:"user_list" query:"user_list"`
}

type FriendUser struct {
	ID            int64   `thrift:"id,1,required" form:"id,required" json:"id,required" query:"id,required"`
	Name          string  `thrift:"name,2,required" form:"name,required" json:"name,required" query:"name,required"`
	FollowCount   *int64  `thrift:"follow_count,3,optional" form:"follow_count" json:"follow_count,omitempty" query:"follow_count"`
	FollowerCount *int64  `thrift:"follower_count,4,optional" form:"follower_count" json:"follower_count,omitempty" query:"follower_count"`
	IsFollow      bool    `thrift:"is_follow,5,required" form:"is_follow,required" json:"is_follow,required" query:"is_follow,required"`
	Avatar        string  `thrift:"avatar,6,required" form:"avatar,required" json:"avatar,required" query:"avatar,required"`
	Message       *string `thrift:"message,7,optional" form:"message" json:"message,omitempty" query:"message"`
	MsgType       int8    `thrift:"msgType,8,required" form:"msgType,required" json:"msgType,required" query:"msgType,required"`
}

type DouyinRelationFriendListResponse struct {
	StatusCode int64         `thrift:"status_code,1,required" form:"status_code,required" json:"status_code,required" query:"status_code,required"`
	StatusMsg  *string       `thrift:"status_msg,2,optional" form:"status_msg" json:"status_msg,omitempty" query:"status_msg"`
	UserList   []*FriendUser `thrift:"user_list,3" form:"user_list" json:"user_list" query:"user_list"`
}

type DouyinMessageChatResponse struct {
	StatusCode  int64      `thrift:"status_code,1,required" form:"status_code,required" json:"status_code,required" query:"status_code,required"`
	StatusMsg   *string    `thrift:"status_msg,2,optional" form:"status_msg" json:"status_msg,omitempty" query:"status_msg"`
	MessageList []*Message `thrift:"message_list,3" form:"message_list" json:"message_list" query:"message_list"`
}

type Message struct {
	ID         int64  `thrift:"id,1,required" form:"id,required" json:"id,required" query:"id,required"`
	ToUserID   int64  `thrift:"to_user_id,2,required" form:"to_user_id,required" json:"to_user_id,required" query:"to_user_id,required"`
	FromUserID int64  `thrift:"from_user_id,3,required" form:"from_user_id,required" json:"from_user_id,required" query:"from_user_id,required"`
	Content    string `thrift:"content,4,required" form:"content,required" json:"content,required" query:"content,required"`
	CreateTime *int64 `thrift:"create_time,5,optional" form:"create_time" json:"create_time,omitempty" query:"create_time"`
}