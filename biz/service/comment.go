package service

import (
	"douyin/biz/model/api"
	"douyin/constant"
	"douyin/dal/db"
	"github.com/cloudwego/hertz/pkg/app"
)

// CommentAction impl
func CommentAction(req *api.DouyinCommentActionRequest, c *app.RequestContext) api.DouyinCommentActionResponse {
	var resp api.DouyinCommentActionResponse

	userId := c.GetInt64(constant.IdentityKey)
	if req.ActionType == 1 {
		//publish comment
		con, err := db.CreateCommentByUserIdAndVideoIdAndContent(*req, uint64(userId))
		if err != nil {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
				Comment:    nil,
			}
		}
		con1, err := db.SelectUserByUserID(uint(userId))
		if err != nil {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
				Comment:    nil,
			}
		}
		con2, err := db.SelectAuthorIdByVideoId(req.VideoID)
		if err != nil || con2 == 0 {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
				Comment:    nil,
			}
		}
		con1.IsFollow = db.IsFollow(uint64(userId), con2)

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
	} else if req.ActionType == 2 {
		//delete comment
		//查询此评论是否是本人发送的
		isComment, err := db.IsCommentCreatedByMyself(uint64(userId), *req.CommentID)
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
		con, err := db.DeleteCommentByCommentId(*req.CommentID)
		if err != nil {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
				Comment:    nil,
			}
		}
		con1, err := db.SelectUserByUserID(uint(userId))
		if err != nil {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
				Comment:    nil,
			}
		}
		con2, err := db.SelectAuthorIdByVideoId(req.VideoID)
		if err != nil || con2 == 0 {
			return api.DouyinCommentActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
				Comment:    nil,
			}
		}
		con1.IsFollow = db.IsFollow(uint64(userId), con2)

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

	//not finish token to userid
	//make a fake userid
	userId := c.GetInt64(constant.IdentityKey)
	//if []api.Comment is nil so the database don't have the data of this user
	list, err := db.SelectCommentListByUserId(uint64(userId), uint64(req.VideoID))
	if err != nil {
		return api.DouyinCommentListResponse{
			StatusCode:  int64(api.ErrCode_ParamErrCode),
			CommentList: nil,
		}
	}

	resp.StatusCode = 0
	resp.CommentList = list
	resp.StatusMsg = nil

	return resp
}
