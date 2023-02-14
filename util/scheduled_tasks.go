package util

import (
	"douyin/dal/db"
	"github.com/go-redis/redis"
	"github.com/robfig/cron/v3"
	"strconv"
)

func ScheduledInit() {
	crontab := cron.New(cron.WithSeconds())

	updateRedis := func() {
		//find all video id
		videoList, err := db.SelectVideoList()
		if err != nil {
			panic(err)
		}
		//update favorite count
		for i := 0; i < len(videoList); i++ {
			redisKey := strconv.FormatInt(int64(videoList[i].ID), 10) + "_video" + "_like"
			_, err := db.RDB.Get(redisKey).Result()
			if err == redis.Nil {
				db.RDB.Set(redisKey, videoList[i].FavoriteCount, 0)
			}
			//update database
			count := db.RDB.Get(redisKey).Val()
			countInt64, err := strconv.ParseInt(count, 10, 64)
			if err != nil {
				panic(err)
			}
			_, err = db.UpdateFavoriteCount(uint64(videoList[i].ID), countInt64)
			if err != nil {
				panic(err)
			}
		}
	}

	//each hour execute
	spec := "*/5 * * * * ?"
	_, err := crontab.AddFunc(spec, updateRedis)
	if err != nil {
		return
	}
	crontab.Start()
}
