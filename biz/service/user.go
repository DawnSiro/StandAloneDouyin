package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/errno"
	"douyin/util"
)

func Register(req *api.DouyinUserRegisterRequest) (*api.DouyinUserRegisterResponse, error) {
	users, err := db.SelectUserByName(req.Username)
	if err != nil {
		return nil, err
	}
	if len(users) != 0 {
		return nil, errno.UserAlreadyExistErr
	}

	// 进行加密并存储
	password := util.BcryptHash(req.Password)
	userID, err := db.CreateUser(&db.User{
		Username: req.Username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return &api.DouyinUserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  nil,
		UserID:     userID,
		Token:      "",
	}, err
}

// Login check user info
func Login(req *api.DouyinUserLoginRequest) (userID int64, err error) {
	users, err := db.SelectUserByName(req.Username)
	if err != nil {
		return 0, err
	}
	if len(users) == 0 {
		return 0, errno.AuthorizationFailedErr
	}

	u := users[0]
	if !util.BcryptCheck(req.Password, u.Password) {
		return 0, errno.AuthorizationFailedErr
	}
	return int64(u.ID), nil
}

func GetUserInfo(req *api.DouyinUserRequest) (*api.DouyinUserResponse, error) {
	u, err := db.SelectUserByUserID(uint(req.UserID))
	if err != nil {
		return nil, err
	}

	// pack
	return &api.DouyinUserResponse{
		StatusCode: 0,
		StatusMsg:  nil,
		User:       u,
	}, nil
}
