package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/model"
	"douyin/dal/pack"
	"douyin/dal/rdb"
	"douyin/pkg/errno"
	"douyin/pkg/pulsar"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func SendMessage(fromUserID, toUserID uint64, content string) (*api.DouyinMessageActionResponse, error) {
	isFriend := db.IsFriend(fromUserID, toUserID)
	if !isFriend {
		errNo := errno.UserRequestParameterError
		errNo.ErrMsg = "不能给非好友发消息"
		hlog.Error("service.message.SendMessage err:", errNo.Error())
		return nil, errNo
	}
	err := pulsar.GetMessageMQInstance().CreateMessage(fromUserID, toUserID, content)
	if err != nil {
		hlog.Error("service.message.SendMessage err: failed to publish a message, ", err.Error())
		return nil, err
	}

	return &api.DouyinMessageActionResponse{
		StatusCode: errno.Success.ErrCode,
	}, nil
}

func GetMessageChat(userID, oppositeID uint64, preMsgTime int64) (*api.DouyinMessageChatResponse, error) {
	logTag := "service.message.GetMessageChat err:"
	if userID == oppositeID {
		hlog.Error(logTag, errno.UserRequestParameterError.Error())
		return nil, errno.UserRequestParameterError
	}
	// 查询缓存
	messageJsonList, err := rdb.GetMessageChatList(userID, oppositeID, preMsgTime)
	if err == nil {
		// 命中缓存
		messages := make([]*model.Message, len(messageJsonList))
		for i := 0; i < len(messageJsonList); i++ {
			// 这里每次循环都会生成一个新的结构体
			message := model.Message{}
			err := json.Unmarshal([]byte(messageJsonList[i]), &message)
			if err != nil {
				hlog.Error(logTag, errno.UserRequestParameterError.Error())
				continue
			}
			// 这里直接取地址就可以了，不用再复制一次
			messages[i] = &message
		}

		hlog.Info("messages:", messages)

		// 返回结果
		return &api.DouyinMessageChatResponse{
			StatusCode:  errno.Success.ErrCode,
			MessageList: pack.Messages(messages),
		}, nil
	}

	// 缓存未命中则查询数据库
	messages, err := db.GetMessagesByUserIDAndPreMsgTime(userID, oppositeID, preMsgTime)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}

	// 更新缓存
	err = rdb.LoadMessageChatList(userID, oppositeID, messages)
	if err != nil {
		hlog.Error(logTag, err.Error())
		return nil, err
	}

	return &api.DouyinMessageChatResponse{
		StatusCode:  errno.Success.ErrCode,
		MessageList: pack.Messages(messages),
	}, nil
}
