// Code generated by hertz generator.

package api

import (
	"context"
	"douyin/biz/service"
	"douyin/pkg/constant"
	"errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	api "douyin/biz/model/api"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// SendMessage .
// @router /douyin/message/action/ [POST]
func SendMessage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinMessageActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	hlog.Info(req)

	fromUserID := c.GetUint64(constant.IdentityKey)
	resp := new(api.DouyinMessageActionResponse)
	if req.ActionType == constant.SendMessageAction {
		resp, err = service.SendMessage(fromUserID, uint64(req.ToUserID), req.Content)
	} else {
		err = errors.New("action type error")
	}

	if err != nil {
		c.JSON(consts.StatusOK, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// GetMessageChat .
// @router /douyin/message/chat/ [GET]
func GetMessageChat(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinMessageChatRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Info(err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	hlog.Info(req)

	userID := c.GetUint64(constant.IdentityKey)
	resp, err := service.GetMessageChat(userID, uint64(req.ToUserID))
	if err != nil {
		c.JSON(consts.StatusOK, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}
