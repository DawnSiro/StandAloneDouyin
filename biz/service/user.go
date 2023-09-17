package service

import (
	"strconv"

	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/model"
	"douyin/dal/pack"
	"douyin/dal/rdb"
	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"douyin/pkg/util"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Register(username, password string) (*api.DouyinUserRegisterResponse, error) {
	user, err := db.SelectUserByName(username)
	if err != nil {
		hlog.Error("service.user.Register err:", err.Error())
		return nil, err
	}
	if user.ID != 0 {
		hlog.Error("service.user.Register err:", errno.UsernameAlreadyExistsError.Error())
		return nil, errno.UsernameAlreadyExistsError
	}

	// 进行加密并存储
	encryptedPassword := util.BcryptHash(password)
	userID, err := db.CreateUser(&model.User{
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

	// 将 UserID 添加到布隆过滤器中
	global.UserIDBloomFilter.AddString(strconv.FormatUint(userID, 10))

	// TODO 预热点赞视频列表逻辑
	// 0 关注用户加一个数来维持 redis key 的存在
	err = rdb.SetFollowUserIDSet(userID, []uint64{0})
	if err != nil {
		hlog.Error("service.user.Register err:", err.Error())
	}

	return &api.DouyinUserRegisterResponse{
		StatusCode: errno.Success.ErrCode,
		UserID:     int64(userID),
		Token:      token,
	}, nil
}

func Login(username, password string) (*api.DouyinUserLoginResponse, error) {
	usernameLoginKey := constant.LoginFailCounterRedisPrefix + username

	usernameLogin, err := global.UserRC.Get(usernameLoginKey).Result()
	var usernameLoginInt int
	if err == nil {
		usernameLoginInt, _ = strconv.Atoi(usernameLogin)
		if usernameLoginInt >= constant.UserLoginLimit {
			msg := "用户登录次数过多，请稍后重试"
			return &api.DouyinUserLoginResponse{StatusCode: errno.UserLoginError.ErrCode, StatusMsg: &msg}, nil
		}
	} else {
		global.UserRC.Set(usernameLoginKey, "0", constant.UserLoginLimitTime)
	}

	user, err := db.SelectUserByName(username)
	if err != nil {
		// 记录登录次数
		global.UserRC.Set(usernameLoginKey, usernameLoginInt+1, constant.UserLoginLimitTime)
		hlog.Error("service.user.Login err:", err.Error())
		return nil, err
	}
	if user.ID == 0 {
		// 记录登录次数
		global.UserRC.Set(usernameLoginKey, usernameLoginInt+1, constant.UserLoginLimitTime)
		hlog.Error("service.user.Login err:", errno.UserAccountDoesNotExistError.Error())
		return nil, errno.UserAccountDoesNotExistError
	}

	if !util.BcryptCheck(password, user.Password) {
		// 密码错误记录登录次数
		global.UserRC.Set(usernameLoginKey, usernameLoginInt+1, constant.UserLoginLimitTime)
		hlog.Error("service.user.Login err:", errno.UserPasswordError.Error())
		return nil, errno.UserPasswordError
	}
	token, err := util.SignToken(user.ID)

	// TODO 预热缓存
	go func() {
		// TODO 将点赞视频ID列表
		// TODO 将关注列表用户ID 的 Set 集合添加到缓存中
		followUserIDSet, err := db.SelectFollowUserIDSet(user.ID)
		if err != nil {
			hlog.Error("service.user.Login err:", err.Error())
			return
		}
		// 0 关注用户加一个数来维持 redis key 的存在
		if len(followUserIDSet) == 0 {
			followUserIDSet = []uint64{0}
		}
		err = rdb.SetFollowUserIDSet(user.ID, followUserIDSet)
		if err != nil {
			hlog.Error("service.user.Login err:", err.Error())
		}
	}()

	return &api.DouyinUserLoginResponse{
		StatusCode: errno.Success.ErrCode,
		UserID:     int64(user.ID),
		Token:      token,
	}, nil
}

// GetUserInfo 获取用户主页信息
// userID 为获取者
// infoUserID 为主页的用户
func GetUserInfo(userID, infoUserID uint64) (*api.DouyinUserResponse, error) {
	// 使用布隆过滤器判断用户ID是否存在
	if !global.UserIDBloomFilter.TestString(strconv.FormatUint(infoUserID, 10)) {
		hlog.Error("service.comment.GetCommentList err: 布隆过滤器拦截")
		return nil, errno.UserRequestParameterError
	}

	// 查询用户信息
	u, err := rdb.GetUserInfo(infoUserID)
	// 缓存未命中则查询数据库
	if err != nil {
		hlog.Error("service.user.GetUserInfo err:", err.Error())
		u, err = db.SelectUserByID(infoUserID)
		if err != nil {
			hlog.Error("service.user.GetUserInfo err:", err.Error())
			return nil, err
		}
		// 设置缓存
		err = rdb.SetUserInfo(u)
		if err != nil {
			hlog.Error("service.user.GetUserInfo err:", err.Error())
		}
	}

	// 判断是否关注
	var isFollow bool
	// 已登录的用户才需要查询是否关注
	if userID != 0 {
		isFollow, err = rdb.IsFollow(userID, infoUserID)
		// 缓存不存在就查数据库
		if err != nil {
			isFollow = db.IsFollow(userID, infoUserID)
			// 关注用户ID部分数据异步更新
			go func() {
				set, err := db.SelectFollowUserIDSet(userID)
				if err != nil {
					hlog.Error("service.user.GetUserInfo err:", err.Error())
					return
				}
				// 0赞用户加入占位符保证缓存存在
				if len(set) == 0 {
					set = []uint64{0}
				}
				err = rdb.SetFollowUserIDSet(userID, set)
				if err != nil {
					hlog.Error("service.user.GetUserInfo err:", err.Error())
				}
			}()
		}
	}

	// 更新缓存过期时间
	err = rdb.ExpireUserInfo(infoUserID)
	if err != nil {
		hlog.Error("service.user.GetUserInfo err:", err.Error())
	}

	return &api.DouyinUserResponse{
		StatusCode: errno.Success.ErrCode,
		User:       pack.UserInfo(u, isFollow),
	}, nil
}
