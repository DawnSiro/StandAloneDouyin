package pulsar

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"douyin/dal/db"
	"douyin/pkg/constant"
	"douyin/pkg/global"
)

type FollowActionMessage struct {
	UpID   uint64
	FanID  uint64
	Action int
}

type FollowActionMQ struct {
	*PulsarMQ
}

var (
	instance *FollowActionMQ
	once     sync.Once
)

func GetFollowActionMQInstance() *FollowActionMQ {
	// 懒汉式单例模式，同时保证线程安全
	if instance == nil {
		once.Do(func() {
			instance = newFollowActionMQ()
		})
	}
	return instance
}

func newFollowActionMQ() *FollowActionMQ {
	res := &FollowActionMQ{
		PulsarMQ: NewPulsarMQ(global.PulsarClient, constant.FollowActionTopic, constant.FollowActionSubscription),
	}
	res.RunConsume()
	return res
}

func (mq *FollowActionMQ) Close() {
	mq.Producer.Close()
	mq.Consumer.Close()
}

func (mq *FollowActionMQ) Consume() error {
	hlog.Info("follow action consumer start")
	for {
		msg, err := mq.Consumer.Receive(context.Background())
		if err != nil {
			return err
		}
		hlog.Debugf("follow action consumer: recieve message (id=%v)", msg.ID())

		var res FollowActionMessage
		err = json.Unmarshal(msg.Payload(), &res)
		if err != nil {
			// 解析错误后丢弃信息但不终止
			hlog.Errorf("follow action consumer: parse message failed (id=%v)", msg.ID())
			err = mq.Consumer.Ack(msg)
			if err != nil {
				return err
			}
			continue
		}

		err = mq.Consumer.Ack(msg)
		if err != nil {
			return err
		}
		hlog.Debugf("follow action consumer: acknowlege message (id=%v)", msg.ID())

		switch res.Action {
		case 1:
			db.Follow(res.FanID, res.UpID)
		case 2:
			db.CancelFollow(res.FanID, res.UpID)
		}
		if err != nil {
			hlog.Errorf("follow action consumer: db error: %v, message (id=%v)", err, msg.ID()) // 数据库错误打印日志，但不停止逻辑
		} else {
			hlog.Debugf("follow action consumer: handle a message successfully")
		}
	}
}

func (mq *FollowActionMQ) RunConsume() {
	go func() {
		err := mq.Consume()
		if err != nil {
			hlog.Errorf("follow action cosumer error: %v", err)
		}
	}()
}

func (mq *FollowActionMQ) FollowAction(upid, fanid uint64) error {
	msg := FollowActionMessage{upid, fanid, constant.Follow}
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = mq.Producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: payload,
	})
	return err
}

func (mq *FollowActionMQ) CancelFollowAction(upid, fanid uint64) error {
	msg := FollowActionMessage{upid, fanid, constant.CancelFollow}
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = mq.Producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: payload,
	})
	return err
}
