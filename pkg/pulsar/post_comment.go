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

type PostCommentMessage db.Comment

type PostCommentMQ struct {
	*PulsarMQ
}

var (
	pcmq *PostCommentMQ
	pcOnce sync.Once
)

func GetPostCommentMQInstance() *PostCommentMQ{
	// 懒汉式单例模式，同时保证线程安全
	if pcmq == nil {
		pcOnce.Do(func() {
			pcmq = newPostCommentMQ()
		})
	}
	return pcmq
}

// 私有化创建实例函数
func newPostCommentMQ() *PostCommentMQ{
	res := &PostCommentMQ{
		PulsarMQ: NewPulsarMQ(global.PulsarClient, constant.CommentActionTopic, constant.CommentActionSubscription),
	}
	res.RunConsume(res.Consume)
	return res
}

func (mq *PostCommentMQ) Consume() error {
	hlog.Info("post comment consumer start")
	for {
		msg, err := mq.Consumer.Receive(context.Background())
		if err != nil {
			return err
		}
		hlog.Debugf("post comment consumer: recieve message (id=%v)", msg.ID())

		err = mq.Consumer.Ack(msg)
		if err != nil {
			return err
		}
		hlog.Debugf("post comment consumer: acknowlege message (id=%v)", msg.ID())

		var res db.Comment
		err = json.Unmarshal(msg.Payload(), &res)
		if err != nil {
			// 解析错误后丢弃信息但不终止
			hlog.Errorf("post comment consumer: parse message failed (id=%v)", msg.ID())
			// TODO: delete data in redis
			continue
		}

		_, err = db.CreateComment(&res)
		if err != nil {
			hlog.Errorf("post comment consumer: db error: %v, message (id=%v)", err, msg.ID()) // 数据库错误打印日志，但不停止逻辑
			// TODO: delete data in redis
		} else {
			hlog.Debugf("post comment consumer: handle a message successfully")
		}
	}
}

func (mq *PostCommentMQ) PostComment (msg PostCommentMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = mq.Producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: payload,
	})
	return err
}
