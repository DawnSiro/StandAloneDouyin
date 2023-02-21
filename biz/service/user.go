package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/errno"
	"douyin/pkg/util"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Register(username, password string) error {
	users, err := db.SelectUserByName(username)
	if err != nil {
		hlog.Error("service.user.Register err:", err.Error())
		return err
	}
	if len(users) != 0 {
		hlog.Error("service.user.Register err:", errno.UsernameAlreadyExistsError.Error())
		return errno.UsernameAlreadyExistsError
	}

	// 进行加密并存储
	encryptedPassword := util.BcryptHash(password)
	_, err = db.CreateUser(&db.User{
		Username: username,
		Password: encryptedPassword,
	})
	if err != nil {
		hlog.Error("service.user.Register err:", err.Error())
		return err
	}
	return nil
}

// Login check user info
func Login(username, password string) (userID uint64, err error) {
	users, err := db.SelectUserByName(username)
	if err != nil {
		hlog.Error("service.user.Login err:", err.Error())
		return 0, err
	}
	if len(users) == 0 {
		hlog.Error("service.user.Login err:", errno.UserAccountDoesNotExistError.Error())
		return 0, errno.UserAccountDoesNotExistError
	}

	u := users[0]
	if !util.BcryptCheck(password, u.Password) {
		hlog.Error("service.user.Login err:", errno.UserPasswordError.Error())
		return 0, errno.UserPasswordError
	}
	return u.ID, nil
}

func GetUserInfo(userID, infoUserID uint64) (*api.DouyinUserResponse, error) {
	u, err := db.SelectUserByID(infoUserID)
	if err != nil {
		hlog.Error("service.user.GetUserInfo err:", err.Error())
		return nil, err
	}

	// pack
	userInfo := pack.UserInfo(u, db.IsFollow(userID, infoUserID))
	return &api.DouyinUserResponse{
		StatusCode: errno.Success.ErrCode,
		User:       userInfo,
	}, nil
}
