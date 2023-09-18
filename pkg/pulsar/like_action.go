package pulsar

import (
	"context"
	"sync"

	"douyin/dal/db"
	"douyin/pkg/constant"
	"douyin/pkg/global"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
)

type LikeActionMessage struct {
	UserID  uint64
	VideoID uint64
	Action  int
}

type LikeActionMQ struct {
	*MQ
}

var (
	lmq   *LikeActionMQ
	lOnce sync.Once
)

func GetLikeActionMQInstance() *LikeActionMQ {
	// 懒汉式单例模式，同时保证线程安全
	lOnce.Do(func() {
		lmq = newLikeActionMQ()
	})
	return lmq
}

// 私有化创建实例函数
func newLikeActionMQ() *LikeActionMQ {
	res := &LikeActionMQ{
		MQ: NewPulsarMQ(global.PulsarClient, constant.LikeActionTopic, constant.LikeActionSubscription),
	}
	res.RunConsume(res.Consume)
	return res
}

func (mq *LikeActionMQ) Consume() error {
	hlog.Info("like action consumer start")
	for {
		msg, err := mq.Consumer.Receive(context.Background())
		if err != nil {
			return err
		}
		hlog.Debugf("like action consumer: receive message (id=%v)", msg.ID())

		err = mq.Consumer.Ack(msg)
		if err != nil {
			return err
		}
		hlog.Debugf("like action consumer: acknowledge message (id=%v)", msg.ID())

		var res LikeActionMessage
		err = json.Unmarshal(msg.Payload(), &res)
		if err != nil {
			// 解析错误后丢弃信息但不终止
			hlog.Errorf("like action consumer: parse message failed (id=%v)", msg.ID())
			// TODO: delete data in redis
			continue
		}

		switch res.Action {
		case constant.Favorite:
			err = db.FavoriteVideo(res.UserID, res.VideoID)
		case constant.CancelFavorite:
			err = db.CancelFavoriteVideo(res.UserID, res.VideoID)
		}
		if err != nil {
			hlog.Errorf("like action consumer: db error: %v, message (id=%v)", err, msg.ID()) // 数据库错误打印日志，但不停止逻辑
			// TODO: delete data in redis
		} else {
			hlog.Debugf("like action consumer: handle a message successfully")
		}
	}
}

func (mq *LikeActionMQ) LikeAction(userID, videoID uint64) error {
	msg := LikeActionMessage{userID, videoID, constant.Favorite}
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = mq.Producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: payload,
	})
	return err
}

func (mq *LikeActionMQ) CancelLikeAction(userID, videoID uint64) error {
	msg := LikeActionMessage{userID, videoID, constant.CancelFavorite}
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = mq.Producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: payload,
	})
	return err
}
