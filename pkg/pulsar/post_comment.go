package pulsar

import (
	"context"
	"douyin/dal/model"
	"encoding/json"
	"sync"

	"douyin/dal/db"
	"douyin/pkg/constant"
	"douyin/pkg/global"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type PostCommentMessage model.Comment

type PostCommentMQ struct {
	*MQ
}

var (
	pCMQ   *PostCommentMQ
	pcOnce sync.Once
)

func GetPostCommentMQInstance() *PostCommentMQ {
	// 懒汉式单例模式，同时保证线程安全
	pcOnce.Do(func() {
		pCMQ = newPostCommentMQ()
	})
	return pCMQ
}

// 私有化创建实例函数
func newPostCommentMQ() *PostCommentMQ {
	res := &PostCommentMQ{
		MQ: NewPulsarMQ(global.PulsarClient, constant.CommentActionTopic, constant.CommentActionSubscription),
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
		hlog.Debugf("post comment consumer: receive message (id=%v)", msg.ID())

		err = mq.Consumer.Ack(msg)
		if err != nil {
			return err
		}
		hlog.Debugf("post comment consumer: acknowledge message (id=%v)", msg.ID())

		var res model.Comment
		err = json.Unmarshal(msg.Payload(), &res)
		if err != nil {
			// 解析错误后丢弃信息但不终止
			hlog.Errorf("post comment consumer: parse message failed (id=%v)", msg.ID())
			continue
		}

		_, err = db.CreateComment(&res)
		if err != nil {
			hlog.Errorf("post comment consumer: db error: %v, message (id=%v)", err, msg.ID()) // 数据库错误打印日志，但不停止逻辑
		} else {
			hlog.Debugf("post comment consumer: handle a message successfully")
		}
	}
}

func (mq *PostCommentMQ) PostComment(msg PostCommentMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = mq.Producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: payload,
	})
	return err
}
