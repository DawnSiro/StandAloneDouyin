package pulsar

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type MQ struct {
	Topic        string
	Subscription string
	Producer     pulsar.Producer
	Consumer     pulsar.Consumer
}

func NewPulsarMQ(client pulsar.Client, topic, subscription string) *MQ {
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		hlog.Fatalf("Failed to create producer: %v", err) // 创建失败将影响业务正常性，因此直接终止程序
	}

	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: subscription,
	})
	if err != nil {
		hlog.Fatalf("Failed to create consumer: %v", err)
	}

	return &MQ{producer.Topic(), consumer.Subscription(), producer, consumer}
}

func (mq *MQ) RunConsume(f func() error) {
	go func() {
		err := f()
		if err != nil {
			hlog.Errorf("follow action consumer error: %v", err)
		}
	}()
}

func (mq *MQ) Close() {
	mq.Producer.Close()
	mq.Consumer.Close()
}
