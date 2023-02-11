// Code generated by hertz generator.

package api

import (
	"context"
	"douyin/biz/service"

	api "douyin/biz/model/api"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// FavoriteVideo .
// @router /douyin/favorite/action/ [POST]
func FavoriteVideo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinFavoriteActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.DouyinFavoriteActionResponse)

	c.JSON(consts.StatusOK, resp)
}

// GetFavoriteList .
// @router /douyin/favorite/list/ [GET]
func GetFavoriteList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinFavoriteListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp, err := service.FavoriteList(&req)
	if err != nil {
		//TODO: 有问题
		c.JSON(consts.StatusOK, err)
	}

	c.JSON(consts.StatusOK, resp)
}
