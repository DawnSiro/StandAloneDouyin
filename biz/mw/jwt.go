package mw

import (
	"context"
	"douyin/biz/model/api"
	"douyin/biz/service"
	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"time"

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
			// 这里的返回值可以通过 c.Get() 或者 c.GetInt64() 去取到
			return uint64(claims[constant.IdentityKey].(float64))
		},
		// 用于设置登陆成功后为向 token 中添加自定义负载信息的函数
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			// 注意这里的数据类型要和下面的 LoginResponse 函数中的类型一致
			if v, ok := data.(uint64); ok {
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
			userID, err := service.Login(req.Username, req.Password)
			if err != nil {
				return 0, err
			}
			hlog.Info(userID)
			// 设置 userID 到请求上下文
			c.Set(constant.IdentityKey, userID)
			return userID, nil
		},
		// 用于设置登录的响应函数
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			// 可以通过 Get 去取放在 请求上下文 的数据
			userID := c.GetUint64(constant.IdentityKey)
			hlog.Info(userID, " ", token)
			c.JSON(http.StatusOK, api.DouyinUserLoginResponse{
				StatusCode: errno.Success.ErrCode,
				UserID:     int64(userID),
				Token:      token,
			})
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(http.StatusOK, utils.H{
				"status_code": errno.UserIdentityVerificationFailedError.ErrCode,
				"status_msg":  errno.UserIdentityVerificationFailedError.Error(),
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
