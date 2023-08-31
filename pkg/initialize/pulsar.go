package initialize

import (
	"log"

	"github.com/apache/pulsar-client-go/pulsar"

	"douyin/pkg/global"
)

func Pulsar() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{})
	if err != nil {
		// TODO: handle error
		log.Fatal(err)
	}
	global.PulsarClient = client
}
