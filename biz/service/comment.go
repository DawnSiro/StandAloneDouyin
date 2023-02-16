package service

import (
	"douyin/biz/model/api"
	"douyin/constant"
	"douyin/dal/db"
	"douyin/util"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/go-redis/redis"
	"strconv"
)

// CommentAction impl
func CommentAction(req *api.DouyinCommentActionRequest, c *app.RequestContext) api.DouyinCommentActionResponse {
	var resp api.DouyinCommentActionResponse
	filterErr := "评论带有敏感词"

	userID := c.GetInt64(constant.IdentityKey)

	//删除redis评论列表缓存
	commentListKey := strconv.FormatInt(req.VideoID, 10) + "_video_" + strconv.FormatInt(userID, 10) + "_userId" + "_comments"
	db.RDB.Del(commentListKey)

	if req.ActionType == constant.PostComment {
		//publish comment
		//检测是否带有敏感词
		if util.IsWordsFilter(*req.CommentText) {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_ParamErrCode),
				StatusMsg:  &filterErr,
			}
		}
		con, err := db.CreateComment(uint64(req.VideoID), *req.CommentText, uint64(userID))
		if err != nil {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
				Comment:    nil,
			}
		}
		con1, err := db.SelectUserByUserID(uint(userID))
		if err != nil {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
				Comment:    nil,
			}
		}
		con2, err := db.SelectAuthorIDByVideoID(req.VideoID)
		if err != nil || con2 == 0 {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
				Comment:    nil,
			}
		}
		con1.IsFollow = db.IsFollow(uint64(userID), con2)

		//update video.comment_count
		_, err = db.IncreaseFavoriteCount(uint64(req.VideoID))

		resp.StatusCode = 0
		resp.StatusMsg = new(string)
		resp.Comment = &api.Comment{
			ID:         int64(con.ID),
			User:       con1,
			Content:    con.Content,
			CreateDate: con.CreatedAt.String(),
		}
	} else if req.ActionType == constant.DeleteComment {
		//delete comment
		//查询此评论是否是本人发送的
		isComment, err := db.IsCommentCreatedByMyself(uint64(userID), *req.CommentID)
		if err != nil {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
				Comment:    nil,
			}
		}
		//非本人评论
		if !isComment {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_AuthorizationFailedErrCode),
				Comment:    nil,
			}
		}
		con, err := db.DeleteCommentByCommentID(*req.CommentID)
		if err != nil {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
				Comment:    nil,
			}
		}
		con1, err := db.SelectUserByUserID(uint(userID))
		if err != nil {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
				Comment:    nil,
			}
		}
		con2, err := db.SelectAuthorIDByVideoID(req.VideoID)
		if err != nil || con2 == 0 {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
				Comment:    nil,
			}
		}
		con1.IsFollow = db.IsFollow(uint64(userID), con2)

		//update video.comment_count
		_, err = db.ReduceFavoriteCount(uint64(req.VideoID))

		resp.StatusCode = 0
		resp.StatusMsg = new(string)
		resp.Comment = &api.Comment{
			ID:         int64(con.ID),
			User:       con1,
			Content:    con.Content,
			CreateDate: con.CreatedAt.String(),
		}
	} else {
		//type is illegal
		return api.DouyinCommentActionResponse{
			StatusCode: int64(api.ErrCode_ParamErrCode),
			Comment:    new(api.Comment),
		}
	}

	return resp
}

func CommentList(req *api.DouyinCommentListRequest, c *app.RequestContext) api.DouyinCommentListResponse {
	var resp api.DouyinCommentListResponse

	userID := c.GetInt64(constant.IdentityKey)

	commentListKey := strconv.FormatInt(req.VideoID, 10) + "_video_" + strconv.FormatInt(userID, 10) + "_userId" + "_comments"
	commentList, err := db.RDB.Get(commentListKey).Result()
	if err == redis.Nil {
		//find like_count in mysql
		list, err := db.SelectCommentListByUserID(uint64(userID), uint64(req.VideoID))
		if err != nil {
			return api.DouyinCommentListResponse{
				StatusCode:  int64(api.ErrCode_ParamErrCode),
				CommentList: nil,
			}
		}
		//序列化
		marshalList, _ := json.Marshal(list)
		db.RDB.Set(commentListKey, marshalList, 0)
	}

	commentList, _ = db.RDB.Get(commentListKey).Result()
	//反序列化
	var list []*api.Comment
	err = json.Unmarshal([]byte(commentList), &list)
	if err != nil {
		err := errors.New("unmarshal error")
		if err != nil {
			return api.DouyinCommentListResponse{}
		}
	}

	resp.StatusCode = 0
	resp.CommentList = list
	resp.StatusMsg = nil

	return resp
}
