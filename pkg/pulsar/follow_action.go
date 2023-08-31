package pulsar

import (
	"context"
	"encoding/json"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"douyin/dal/db"
	"douyin/pkg/constant"
)

type FollowActionMessage struct {
	UpID   uint64
	FanID  uint64
	Action int
}

type FollowActionMQ struct {
	Topic        string
	Subscription string
	Producer     pulsar.Producer
	Consumer     pulsar.Consumer
}

func NewFollowActionMQ(client pulsar.Client) (*FollowActionMQ, error) {
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: constant.FollowActionTopic,
	})
	if err != nil {
		return nil, err
	}

	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            constant.FollowActionTopic,
		SubscriptionName: constant.FollowActionTopic + "sub",
	})
	if err != nil {
		return nil, err
	}

	return &FollowActionMQ{constant.FollowActionTopic, consumer.Subscription(), producer, consumer}, nil
}

func (mq *FollowActionMQ) Close() {
	mq.Producer.Close()
	mq.Consumer.Close()
}

func (mq *FollowActionMQ) Consume() error {
	hlog.Info("service.relation.Follow: follow action consumer start")
	for {
		msg, err := mq.Consumer.Receive(context.Background())
		if err != nil {
			return err
		}
		hlog.Debugf("Recieve message(id=%v)", msg.ID().String())

		var res FollowActionMessage
		err = json.Unmarshal(msg.Payload(), &res)
		if err != nil {
			return err
		}

		err = mq.Consumer.Ack(msg)
		if err != nil {
			return err
		}
		hlog.Debugf("Acknowlege message(id=%v)", msg.ID().String())

		switch res.Action {
		case 1:
			db.Follow(res.FanID, res.UpID)
		case 2:
			db.CancelFollow(res.FanID, res.UpID)
		}
		if err != nil {
			hlog.Errorf("follow action cosumer db error: %v, message id: %v", err, msg.ID()) // 数据库错误打印日志，但不停止逻辑
		}
	}
}

func (mq *FollowActionMQ) RunConsume() {
	go func() {
		err := mq.Consume()
		if err != nil {
			hlog.Errorf("follow action consumer error: %v", err)
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
