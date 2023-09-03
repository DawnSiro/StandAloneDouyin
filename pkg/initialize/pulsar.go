package initialize

import (
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"douyin/pkg/global"
)

func Pulsar() {
	client, err := client()
	if err != nil {
		hlog.Fatal(err)  // 失败后直接终止程序
	}

	global.PulsarClient = client
	hlog.Info("pulsar initialized successfullly")
}

func client() (client pulsar.Client, err error) {
	client, err = pulsar.NewClient(pulsar.ClientOptions{
		URL:               global.Config.PulsarConfig.URL,
		ConnectionTimeout: time.Second * time.Duration(global.Config.PulsarConfig.ConnectionTimeout),
		OperationTimeout:  time.Second * time.Duration(global.Config.PulsarConfig.OperationTimeout),
		// TODO: more config
	})
	if err != nil {
		return
	}

	// 检验连接是否正常
	_, err = client.CreateProducer(pulsar.ProducerOptions{
		Topic: "ping",
	})
	return
}
