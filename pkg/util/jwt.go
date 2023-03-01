package util

import (
	"time"

	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"douyin/pkg/global"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	ID uint64
	jwt.RegisteredClaims
}

// SignToken 签发 Token
func SignToken(userID uint64) (string, error) {
	// 将签名字符串转化为 byte 切片
	signingKey := []byte(global.Config.JWTConfig.SigningKey)
	// 配置 userClaims ,并生成 token
	claims := UserClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constant.TokenTimeOut)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

// ParseToken 解析 token
func ParseToken(tokenString string) (*UserClaims, error) {
	// 将签名字符串转化为 byte 切片
	signingKey := []byte(global.Config.JWTConfig.SigningKey)
	// 解析 token 信息
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		hlog.Error("util.jwt.ParseToken err:", err.Error())
		return nil, err
	} else if token == nil {
		hlog.Error("util.jwt.ParseToken err:", err.Error())
		return nil, errno.UserIdentityVerificationFailedError
	}
	if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	}
	hlog.Error("util.jwt.ParseToken err:", err.Error())
	return nil, errno.UserIdentityVerificationFailedError
}
