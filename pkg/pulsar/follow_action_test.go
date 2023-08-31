package pulsar
import (
	"testing"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
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

func TestFollowActionMQProducer(t *testing.T) {
	client, err := GetClient()
	if err != nil {
		t.Fatal(err)
	}

	fmq, err := NewFollowActionMQ(client)
	if err != nil {
		t.Fatal(err)
	}

	// fmq.RunConsume()

	err = fmq.FollowAction(1, 2)
	if err != nil {
		t.Error(err)
	}

	err = fmq.CancelFollowAction(1, 2)
	if err != nil {
		t.Error(err)
	}
}

// func TestFollowActionConsumer(t *testing.T) {
// 	client, err := GetClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	fmq, err := NewFollowActionMQ(client)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	go func () {
// 		err := fmq.FollowAction(1, 2)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		err = fmq.CancelFollowAction(1, 2)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}()

// 	err = fmq.Consume()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
