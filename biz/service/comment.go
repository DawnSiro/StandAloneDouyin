package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
)

// CommentAction impl
func CommentAction(req *api.DouyinCommentActionRequest) api.DouyinCommentActionResponse {
	var resp api.DouyinCommentActionResponse
	serverError := "服务器内部错误"

	if req.ActionType == 1 {
		//publish comment
		con, err := db.CreateCommentByUserIdAndVideoIdAndContent(*req)
		if err != nil {
			return api.DouyinCommentActionResponse{
				StatusCode: 1001,
				StatusMsg:  &serverError,
				Comment:    nil,
			}
		}
		//TODO:miss Avatar
		con1, err := db.SelectUserByUserId(uint(req.UserID))
		if err != nil {
			return api.DouyinCommentActionResponse{
				StatusCode: 1001,
				StatusMsg:  &serverError,
				Comment:    nil,
			}
		}
		con2, err := db.SelectAuthorIdByVideoId(req.VideoID)
		if err != nil || con2 == 0 {
			return api.DouyinCommentActionResponse{
				StatusCode: 1001,
				StatusMsg:  &serverError,
				Comment:    nil,
			}
		}
		con1.IsFollow = db.IsFollow(uint64(req.UserID), con2)

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
		con, err := db.DeleteCommentByCommentId(*req.CommentID)
		if err != nil {
			return api.DouyinCommentActionResponse{
				StatusCode: 1001,
				StatusMsg:  &serverError,
				Comment:    nil,
			}
		}
		//TODO:miss Avatar
		con1, err := db.SelectUserByUserId(uint(req.UserID))
		if err != nil {
			return api.DouyinCommentActionResponse{
				StatusCode: 1001,
				StatusMsg:  &serverError,
				Comment:    nil,
			}
		}
		con2, err := db.SelectAuthorIdByVideoId(req.VideoID)
		if err != nil || con2 == 0 {
			return api.DouyinCommentActionResponse{
				StatusCode: 1001,
				StatusMsg:  &serverError,
				Comment:    nil,
			}
		}
		con1.IsFollow = db.IsFollow(uint64(req.UserID), con2)

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
			StatusCode: 1002,
			StatusMsg:  &serverError,
			Comment:    new(api.Comment),
		}
	}

	return resp
}

func CommentList(req *api.DouyinCommentListRequest) api.DouyinCommentListResponse {
	var resp api.DouyinCommentListResponse
	serverError := "服务器内部错误"

	//TODO:token -> userId
	//not finish token to userid
	//make a fake userid
	userId := 1
	//if []api.Comment is nil so the database don't have the data of this user
	list, err := db.SelectCommentListByUserId(uint64(userId), uint64(req.VideoID))
	if err != nil {
		return api.DouyinCommentListResponse{
			StatusCode:  1001,
			StatusMsg:   &serverError,
			CommentList: nil,
		}
	}

	resp.StatusCode = 0
	resp.CommentList = list
	resp.StatusMsg = nil

	return resp
}
