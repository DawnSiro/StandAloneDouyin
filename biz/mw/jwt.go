package mw

import (
	"context"
	"douyin/biz/model/api"
	"douyin/biz/service"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"time"

	"douyin/constant"
	"douyin/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"
)

var JwtMiddleware *jwt.HertzJWTMiddleware

func InitJWT() {
	JwtMiddleware, _ = jwt.New(&jwt.HertzJWTMiddleware{
		Key:           []byte(constant.SecretKey),
		TokenLookup:   "header: Authorization, query: token, cookie: jwt, form: token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		Timeout:       time.Hour * 12,
		MaxRefresh:    time.Hour * 3,
		IdentityKey:   constant.IdentityKey,
		// 用于设置获取身份信息的函数，默认与 IdentityKey 配合使用
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			//return &api.User{
			//	ID: int64(claims[constant.IdentityKey].(float64)),
			//}
			// 这里的返回值可以通过 c.Get() 或者 c.GetInt64() 去取到
			hlog.Infof("jwt %d", int64(claims[constant.IdentityKey].(float64)))
			return int64(claims[constant.IdentityKey].(float64))
		},
		// 用于设置登陆成功后为向 token 中添加自定义负载信息的函数
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					constant.IdentityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		// 用于设置登录时认证用户信息的函数（必要配置）
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var err error
			var req api.DouyinUserLoginRequest
			if err = c.BindAndValidate(&req); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			if len(req.Username) == 0 || len(req.Password) == 0 {
				return "", jwt.ErrMissingLoginValues
			}
			userID, err := service.Login(&api.DouyinUserLoginRequest{
				Username: req.Username,
				Password: req.Password,
			})
			if err != nil {
				return 0, err
			}
			c.Set(constant.IdentityKey, userID)
			return userID, nil
		},
		// 用于设置登录的响应函数
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			hlog.Info("loginResponse")
			// 可以通过 Get 去取放在 请求上下文 的数据
			userID := c.GetInt64(constant.IdentityKey)
			c.JSON(http.StatusOK, api.DouyinUserLoginResponse{
				StatusCode: errno.Success.ErrCode,
				StatusMsg:  nil,
				UserID:     userID,
				Token:      token,
			})
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(http.StatusOK, utils.H{
				"status_code": errno.AuthorizationFailedErr.ErrCode,
				"status_msg":  message,
			})
		},
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			switch t := e.(type) {
			case errno.ErrNo:
				return t.ErrMsg
			default:
				return t.Error()
			}
		},
	})
}
