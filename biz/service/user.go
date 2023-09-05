package service

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/errno"
	"douyin/pkg/global"
	"douyin/pkg/util"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
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
		StatusCode: errno.Success.ErrCode,
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
		StatusCode: errno.Success.ErrCode,
		UserID:     int64(u.ID),
		Token:      token,
	}, nil
}

func GetUserInfo(userID, infoUserID uint64) (*api.DouyinUserResponse, error) {
	cacheKey := fmt.Sprintf("userinfo:%d", infoUserID)
	if bloomFilter.TestString(cacheKey) {
		cachedData, err := global.UserInfoRC.Get(cacheKey).Result()
		if err == nil {
			// Cache hit, return cached user info
			var cachedResponse api.DouyinUserResponse
			if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err != nil {
				hlog.Error("service.user.GetUserInfo err: Error decoding cached data, ", err.Error())
			} else {
				return &cachedResponse, nil
			}
		}
	}

	// Create a random duration for cache expiration
	minDuration := 24 * time.Hour
	maxDuration := 48 * time.Hour
	cacheDuration := minDuration + time.Duration(rand.Intn(int(maxDuration-minDuration)))

	// Create a WaitGroup for the cache update operation
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	// Check if another thread is updating the cache
	cacheMutex.Lock()
	existingWaitGroup, exists := cacheStatus[cacheKey]
	if exists {
		cacheMutex.Unlock()
		existingWaitGroup.Wait()
		return GetUserInfo(userID, infoUserID)
	}
	// Set cache status flag to indicate cache update is in progress
	cacheStatus[cacheKey] = waitGroup
	cacheMutex.Unlock()

	// Cache miss, query the database
	u, err := db.SelectUserByID(infoUserID)
	if err != nil {
		hlog.Error("service.user.GetUserInfo err:", err.Error())
		// Release cache status flag to allow other threads to update cache
		cacheMutex.Lock()
		delete(cacheStatus, cacheKey)
		cacheMutex.Unlock()
		return nil, err
	}

	// Pack user info
	userInfo := pack.UserInfo(u, db.IsFollow(userID, infoUserID))
	response := &api.DouyinUserResponse{
		StatusCode: errno.Success.ErrCode,
		User:       userInfo,
	}

	// Store user info in Redis cache with the random expiration time
	responseJSON, _ := json.Marshal(response)
	err = global.UserInfoRC.Set(cacheKey, responseJSON, cacheDuration).Err()
	if err != nil {
		hlog.Error("service.user.GetUserInfo err: Error storing data in cache, ", err.Error())
	}

	// Release cache status flag and signal that cache update is done
	cacheMutex.Lock()
	delete(cacheStatus, cacheKey)
	waitGroup.Done()
	cacheMutex.Unlock()

	return response, nil
}
