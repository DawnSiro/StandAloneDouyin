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
		hlog.Fatal("Failed to initialize sonyflake: ", err)
	}
}

func GetSonyflakeID() (uint64, error) {
	id, err := instance.NextID()
	return id, err
}
