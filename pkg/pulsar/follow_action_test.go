package pulsar

import (
	"douyin/pkg/global"
	"testing"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetClient() (client pulsar.Client, err error) {
	client, err = pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://localhost:6650",
		ConnectionTimeout: 30 * time.Second,
		OperationTimeout:  30 * time.Second,
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

func TestFollowActionMQ(t *testing.T) {
	var err error
	global.PulsarClient, err = GetClient()
	if err != nil {
		t.Fatal(err)
	}
	global.DB, err = gorm.Open(mysql.Open("root:root@tcp(127.0.0.1)/douyin"))
	if err != nil {
		t.Fatal(err)
	}

	err = GetFollowActionMQInstance().FollowAction(1, 2)
	if err != nil {
		t.Error(err)
	}

	err = GetFollowActionMQInstance().CancelFollowAction(1, 2)
	if err != nil {
		t.Error(err)
	}
}
