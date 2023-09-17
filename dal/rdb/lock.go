package rdb

import (
	"context"
	"errors"
	"strconv"
	"time"

	"douyin/pkg/global"
	"douyin/pkg/util"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis"
)

// DistributedLock 分布式锁结构体，每次加锁都会生成一个，并且会
type DistributedLock struct {
	// Time to Live 锁的过期时间
	TTL             time.Duration
	RandomValue     uint64
	Key             string
	TryLockInterval time.Duration
	watchDog        chan bool
}

// tryLock 试图去加锁，仅由 Lock 方法进行调用
func (l *DistributedLock) tryLock(rc *redis.Client) (bool, error) {
	// SetNX Not Exist 不存在才能设置成功
	result, err := rc.SetNX(l.Key, l.RandomValue, l.TTL).Result()
	// err != nil 说明报了错，result = false 代表没有错，并且因为键已经存在所以设置失败，result = true 代表设置成功
	if err != nil {
		return false, err
	}
	// 成功之后需要启动 WatchDog 再返回
	if result {
		go l.startWatchDog(rc)
	}
	// 返回就行
	return result, err
}

// Unlock 进行解锁
func (l *DistributedLock) Unlock(rc *redis.Client) error {
	_, err := rc.EvalSha(global.UnLockLuaScriptHash, []string{l.Key}, l.RandomValue).Result()
	// 不管有没有错，都通知已经解锁
	close(l.watchDog)
	// 然后返回错误
	return err
}

// Lock 加分布式锁，由 Redis 的性质保证即使在微服务场景下，锁也能起作用
func (l *DistributedLock) Lock(ctx context.Context, rc *redis.Client) error {
	// 先试着加锁
	result, err := l.tryLock(rc)
	if result {
		return nil
	}
	// 正常没抢到锁 err 还是 nil
	// 如果 err != nil 了，说明有了其他的问题，此时直接结束就好
	if err != nil {
		hlog.Error("dal.rdb.lock.Lock err:", err.Error())
		return err
	}
	// 加锁失败，不断尝试
	ticker := time.NewTicker(l.TryLockInterval)
	defer ticker.Stop()
	for {
		select {
		// 收到 context 中的超时消息则返回
		case <-ctx.Done():
			return errors.New("锁超时了")
		// 定时尝试重新尝试加锁
		case <-ticker.C:
			lock, err := l.tryLock(rc)
			// 加锁成功则返回
			if lock {
				return nil
			}
			// 正常没抢到锁 err 还是 nil
			// 如果 err != nil 了，说明有了其他的问题，此时直接结束就好
			if err != nil {
				return err
			}
		}
	}
}

// startWatchDog 启动 WatchDog，自动定时延长锁时间，如果接收到锁被释放的信号则结束
func (l *DistributedLock) startWatchDog(rc *redis.Client) {
	ticker := time.NewTicker(l.TTL - 200)
	defer ticker.Stop()
	for {
		select {
		// 定时延长锁的过期时间
		case <-ticker.C:
			_, err := rc.Expire(l.Key, l.TTL).Result()
			// 异常或锁已经不存在则不再续期
			if err != nil {
				hlog.Error("dal.rdb.lock.startWatchDog err:", err.Error())
				return
			}
		// 接收到 watchDog 发送的关闭信号，锁已经释放
		case <-l.watchDog:
			return
		}
	}
}

// NewUserKeyLock 创建一个新的 Lock 结构体对象
func NewUserKeyLock(userID uint64, prefix string) DistributedLock {
	value, _ := util.GetSonyFlakeID()
	timeInterval := time.Duration(100)
	watchDog := make(chan bool)
	return DistributedLock{
		Key:             prefix + strconv.FormatUint(userID, 10),
		RandomValue:     value,
		TTL:             500,
		TryLockInterval: timeInterval,
		watchDog:        watchDog,
	}
}
