// Code generated by hertz generator.

package api

import (
	"context"
	api "douyin/biz/model/api"
	"douyin/biz/service"
	"douyin/pkg/constant"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetFeed .
// @router /douyin/feed/ [GET]
func GetFeed(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinFeedRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint64(constant.IdentityKey)
	resp, err := service.GetFeed(req.LatestTime, userID)
	if err != nil {
		c.JSON(consts.StatusOK, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}
