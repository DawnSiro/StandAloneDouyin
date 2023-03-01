package mw

import (
	"context"
	"douyin/pkg/constant"
	"douyin/pkg/global"
	"douyin/pkg/util"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"net/http"
	"time"

	"douyin/biz/model/api"
	"douyin/biz/service"
	"douyin/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"
)

var JwtMiddleware *jwt.HertzJWTMiddleware

func InitJWT() {
	JwtMiddleware, _ = jwt.New(&jwt.HertzJWTMiddleware{
		Key:           []byte(global.Config.JWTConfig.SigningKey),
		TokenLookup:   "header: Authorization, query: token, cookie: jwt, form: token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		Timeout:       constant.TokenTimeOut,
		MaxRefresh:    constant.TokenMaxRefresh,
		IdentityKey:   global.Config.JWTConfig.IdentityKey,
		// 用于设置获取身份信息的函数，默认与 IdentityKey 配合使用
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			// 这里的返回值可以通过 c.Get() 或者 c.GetInt64() 去取到
			return uint64(claims[global.Config.JWTConfig.IdentityKey].(float64))
		},
		// 用于设置登陆成功后为向 token 中添加自定义负载信息的函数
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(uint64); ok {
				return jwt.MapClaims{
					global.Config.JWTConfig.IdentityKey: v,
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
			resp, err := service.Login(req.Username, req.Password)
			if err != nil {
				return 0, err
			}
			// 设置 userID 到请求上下文
			c.Set(global.Config.JWTConfig.IdentityKey, resp.UserID)
			return resp.UserID, nil
		},
		// 用于设置登录的响应函数
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			// 可以通过 Get 去取放在 请求上下文 的数据
			userID := c.GetUint64(global.Config.JWTConfig.IdentityKey)
			c.JSON(http.StatusOK, api.DouyinUserLoginResponse{
				StatusCode: errno.Success.ErrCode,
				UserID:     int64(userID),
				Token:      token,
			})
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(http.StatusOK, utils.H{
				"status_code": errno.UserIdentityVerificationFailedError.ErrCode,
				"status_msg":  errno.UserIdentityVerificationFailedError.ErrMsg,
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

func JWT() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := c.Query("token")
		if token == "" {
			token = c.PostForm("token")
		}
		if token == "" {
			hlog.Error("mw.jwt.ParseToken err:", errno.UserIdentityVerificationFailedError)
			c.JSON(consts.StatusOK, &api.DouyinResponse{
				StatusCode: errno.UserIdentityVerificationFailedError.ErrCode,
				StatusMsg:  errno.UserIdentityVerificationFailedError.ErrMsg,
			})
			c.Abort()
			return
		}
		claim, err := util.ParseToken(token)
		if err != nil {
			hlog.Error("mw.jwt.ParseToken err:", err.Error())
			c.JSON(consts.StatusOK, &api.DouyinResponse{
				StatusCode: errno.UserIdentityVerificationFailedError.ErrCode,
				StatusMsg:  errno.UserIdentityVerificationFailedError.ErrMsg,
			})
			c.Abort()
			return
		}
		hlog.Info("mw.jwt.ParseToken userID:", claim.ID)
		c.Set(global.Config.JWTConfig.IdentityKey, claim.ID)
		c.Next(ctx)
	}
}

// ParseToken 如果 Token 存在，会试着解析 Token，不存在也会放行。主要用于某些登录和未登录都能使用的接口
func ParseToken() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := c.Query("token")
		if token == "" {
			return
		}
		claim, err := util.ParseToken(token)
		if err != nil {
			hlog.Info("mw.jwt.ParseToken err:", err.Error())
		} else {
			hlog.Info("mw.jwt.ParseToken userID:", claim.ID)
			c.Set(global.Config.JWTConfig.IdentityKey, claim.ID)
		}
		c.Next(ctx)
	}
}
