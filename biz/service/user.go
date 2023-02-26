package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/errno"
	"douyin/pkg/util"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Register(username, password string) (*api.DouyinUserRegisterResponse, error) {
	users, err := db.SelectUserByName(username)
	if err != nil {
		hlog.Error("service.user.Register err:", err.Error())
		return nil, err
	}
	if len(users) != 0 {
		hlog.Error("service.user.Register err:", errno.UsernameAlreadyExistsError.Error())
		return nil, errno.UsernameAlreadyExistsError
	}

	// 进行加密并存储
	encryptedPassword := util.BcryptHash(password)
	userID, err := db.CreateUser(&db.User{
		Username: username,
		Password: encryptedPassword,
	})
	if err != nil {
		hlog.Error("service.user.Register err:", err.Error())
		return nil, err
	}
	token, err := util.SignToken(userID)
	if err != nil {
		hlog.Error("service.user.Register err:", err.Error())
		return nil, err
	}
	return &api.DouyinUserRegisterResponse{
		StatusCode: 0,
		UserID:     int64(userID),
		Token:      token,
	}, nil
}

func Login(username, password string) (*api.DouyinUserLoginResponse, error) {
	users, err := db.SelectUserByName(username)
	if err != nil {
		hlog.Error("service.user.Login err:", err.Error())
		return nil, err
	}
	if len(users) == 0 {
		hlog.Error("service.user.Login err:", errno.UserAccountDoesNotExistError.Error())
		return nil, errno.UserAccountDoesNotExistError
	}

	u := users[0]
	if !util.BcryptCheck(password, u.Password) {
		hlog.Error("service.user.Login err:", errno.UserPasswordError.Error())
		return nil, errno.UserPasswordError
	}
	token, err := util.SignToken(u.ID)
	return &api.DouyinUserLoginResponse{
		StatusCode: 0,
		UserID:     int64(u.ID),
		Token:      token,
	}, nil
}

func GetUserInfo(userID, infoUserID uint64) (*api.DouyinUserResponse, error) {
	u, err := db.SelectUserByID(infoUserID)
	if err != nil {
		hlog.Error("service.user.GetUserInfo err:", err.Error())
		return nil, err
	}

	// TODO 使用 Redis Hash 来对用户数据进行缓存

	// pack
	userInfo := pack.UserInfo(u, db.IsFollow(userID, infoUserID))
	return &api.DouyinUserResponse{
		StatusCode: errno.Success.ErrCode,
		User:       userInfo,
	}, nil
}
