package rdb

import (
	"douyin/dal/model"
	"douyin/pkg/global"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"time"
)

// GetMessageChatList 获取某两个用户之间的消息记录
// Redis key: [userID]:[toUserID] 小的在前面
func GetMessageChatList(userID uint64, toUserID uint64, preMsgTime int64) ([]string, error) {
	messageListKey := getMessageKey(userID, toUserID)
	jsonList, err := global.MessageRC.ZRangeByScore(messageListKey, redis.ZRangeBy{
		Min: strconv.FormatInt(preMsgTime, 10),
	}).Result()
	if err != nil {
		hlog.Error("rdb.message.GetMessageChatList err:", err.Error())
		return nil, err
	}
	return jsonList, nil
}

// LoadMessageChatList 加载消息列表
// Redis key: [userID]:[toUserID] 小的在前面
func LoadMessageChatList(userID uint64, toUserID uint64, messageList []*model.Message) error {
	messageListKey := getMessageKey(userID, toUserID)
	for i := 0; i < len(messageList); i++ {
		messageJson, _ := json.Marshal(messageList[i])
		err := global.MessageRC.ZAdd(messageListKey, redis.Z{
			Score:  float64(messageList[i].CreatedTime.UnixMilli()),
			Member: messageJson,
		}).Err()
		if err != nil {
			hlog.Error("rdb.message.LoadMessageChatList err:", err.Error())
			return err
		}
		err = global.MessageRC.Expire(messageListKey, 2*time.Hour).Err()
		if err != nil {
			hlog.Error("rdb.message.LoadMessageChatList err:", err.Error())
			return err
		}
	}
	return nil
}

// AddMessage 添加 Message 到 Redis 中
func AddMessage(userID uint64, toUserID uint64, message *model.Message) error {
	logTag := "rdb.message.AddMessage err:"
	messageListKey := getMessageKey(userID, toUserID)
	messageJson, err := json.Marshal(message)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return err
	}
	// 加入 ZSet 集合
	err = global.MessageRC.ZAdd(messageListKey, redis.Z{
		Score:  float64(message.CreatedTime.UnixMilli()),
		Member: messageJson,
	}).Err()
	if err != nil {
		hlog.Error(logTag, err.Error())
		return err
	}
	err = global.MessageRC.Expire(messageListKey, 2*time.Hour).Err()
	if err != nil {
		hlog.Error(logTag, err.Error())
		return err
	}
	return nil
}

// DelMessage 从 Redis 中删除消息
func DelMessage(userID uint64, toUserID uint64, message *model.Message) error {
	messageListKey := getMessageKey(userID, toUserID)
	messageJson, err := json.Marshal(message)
	if err != nil {
		hlog.Error("rdb.message.DelMessage err:", err.Error())
		return err
	}
	err = global.MessageRC.ZRem(messageListKey, messageJson).Err()
	if err != nil {
		hlog.Error("rdb.message.DelMessage err:", err.Error())
		return err
	}
	return nil
}

// getMessageKey 获取消息列表的 Redis key
func getMessageKey(userID uint64, toUserID uint64) string {
	var builder strings.Builder
	if userID < toUserID {
		builder.WriteString("chat")
		builder.WriteString(strconv.FormatUint(userID, 10))
		builder.WriteString(":")
		builder.WriteString(strconv.FormatUint(toUserID, 10))
	} else {
		builder.WriteString("chat")
		builder.WriteString(strconv.FormatUint(toUserID, 10))
		builder.WriteString(":")
		builder.WriteString(strconv.FormatUint(userID, 10))
	}
	return builder.String()
}
