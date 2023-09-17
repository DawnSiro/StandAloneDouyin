package pulsar

import (
	"context"
	"encoding/json"
	"sync"

	"douyin/dal/db"
	"douyin/pkg/constant"
	"douyin/pkg/global"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var (
	mmq   *MessageMQ
	mOnce sync.Once
)

type MessageMessage struct {
	UID     uint64
	ToUID   uint64
	Context string
}

type MessageMQ struct {
	*MQ
}

func GetMessageMQInstance() *MessageMQ {
	mOnce.Do(func() {
		mmq = newMessageMQ()
	})
	return mmq
}

func newMessageMQ() *MessageMQ {
	res := &MessageMQ{NewPulsarMQ(global.PulsarClient, constant.MessageTopic, constant.MessageSubscription)}
	res.RunConsume(res.Consume)
	return res
}

func (mq *MessageMQ) Consume() error {
	hlog.Info("message consumer start")
	for {
		msg, err := mq.Consumer.Receive(context.Background())
		if err != nil {
			return err
		}
		hlog.Debugf("message consumer: receive message (id=%v)", msg.ID())

		err = mq.Consumer.Ack(msg)
		if err != nil {
			return err
		}
		hlog.Debugf("message consumer: acknowledge message (id=%v)", msg.ID())

		var res MessageMessage
		err = json.Unmarshal(msg.Payload(), &res)
		if err != nil {
			// 解析错误后丢弃信息但不终止
			hlog.Errorf("message consumer: parse message failed (id=%v)", msg.ID())
			continue
		}

		err = db.CreateMessage(res.UID, res.ToUID, res.Context)
		if err != nil {
			hlog.Errorf("message consumer: db error: %v, message (id=%v)", err, msg.ID()) // 数据库错误打印日志，但不停止逻辑
		} else {
			hlog.Debugf("message consumer: handle a message successfully")
		}
	}
}

func (mq *MessageMQ) CreateMessage(userID, toUserID uint64, comment string) error {
	msg := MessageMessage{userID, toUserID, comment}
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = mq.Producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: payload,
	})
	return err
}
