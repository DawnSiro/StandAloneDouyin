package util

import (
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/sony/sonyflake"
)

var (
	instance *sonyflake.Sonyflake
)

func init() {
	var err error
	instance, err = sonyflake.New(sonyflake.Settings{
		StartTime: time.Now(),
	})
	if err != nil {
		hlog.Fatal("Failed to initialize sony flake: ", err)
	}
}

func GetSonyFlakeID() (uint64, error) {
	return instance.NextID()
}
