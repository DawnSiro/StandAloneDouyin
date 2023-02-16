package service

import (
	"douyin/biz/model/api"
	"douyin/constant"
	"douyin/dal/db"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
)

func Follow(req api.DouyinRelationActionRequest, c *app.RequestContext) (api.DouyinRelationActionResponse, error) {
	//先获取到userID
	userID := c.GetInt64(constant.IdentityKey)
	errorText := "请勿重复操作"
	errorText2 := "不能自己关注自己哦"

	if uint64(userID) == uint64(req.ToUserID) {
		return api.DouyinRelationActionResponse{
			StatusCode: int64(api.ErrCode_ParamErrCode),
			StatusMsg:  &errorText2,
		}, errors.New(strconv.FormatInt(int64(api.ErrCode_ParamErrCode), 10))
	}
	isFollow := db.IsFollow(uint64(userID), uint64(req.ToUserID))
	if !isFollow && req.ActionType == constant.Follow {
		//关注操作
		err := db.AddFollow(uint64(userID), uint64(req.ToUserID))
		if err != nil {
			return api.DouyinRelationActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
			}, nil
		}
	} else if isFollow && req.ActionType == constant.CancelFavorite {
		//取消关注
		err := db.DelFollow(uint64(userID), uint64(req.ToUserID))
		if err != nil {
			return api.DouyinRelationActionResponse{
				StatusCode: int64(api.ErrCode_ServiceErrCode),
			}, nil
		}

	} else {
		return api.DouyinRelationActionResponse{
			StatusCode: int64(api.ErrCode_ParamErrCode),
			StatusMsg:  &errorText,
		}, errors.New(strconv.FormatInt(int64(api.ErrCode_ParamErrCode), 10))
	}
	return api.DouyinRelationActionResponse{
		StatusCode: int64(api.ErrCode_SuccessCode),
	}, nil

}

func GetFollowList(req api.DouyinRelationFollowListRequest) (api.DouyinRelationFollowListResponse, error) {
	resp := api.DouyinRelationFollowListResponse{}

	resultList, err := db.GetFollowList(uint64(req.UserID))
	if err != nil {
		return api.DouyinRelationFollowListResponse{
			StatusCode: int64(api.ErrCode_ServiceErrCode),
		}, err
	}

	resp.StatusCode = int64(api.ErrCode_SuccessCode)
	resp.UserList = resultList
	return resp, nil

}

func GetFollowerList(req api.DouyinRelationFollowerListRequest) (api.DouyinRelationFollowerListResponse, error) {
	resp := api.DouyinRelationFollowerListResponse{}

	resultList, err := db.GetFollowerList(uint64(req.UserID))
	if err != nil {
		return api.DouyinRelationFollowerListResponse{
			StatusCode: int64(api.ErrCode_ServiceErrCode),
		}, err
	}

	resp.StatusCode = int64(api.ErrCode_SuccessCode)
	resp.UserList = resultList
	return resp, nil
}

func GetFriendList(req api.DouyinRelationFriendListRequest) (api.DouyinRelationFriendListResponse, error) {
	resp := api.DouyinRelationFriendListResponse{}

	resultList, err := db.GetFriendList(uint64(req.UserID))
	if err != nil {
		return api.DouyinRelationFriendListResponse{
			StatusCode: int64(api.ErrCode_ServiceErrCode),
		}, err
	}

	resp.StatusCode = int64(api.ErrCode_SuccessCode)
	resp.UserList = resultList
	return resp, nil
}
