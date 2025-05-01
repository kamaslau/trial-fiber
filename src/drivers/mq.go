package drivers

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

// MQClient 连接单例
var MQClient *nats.Conn

// ConnectMQ 建立连接
func ConnectMQ() {
	MQ_URL := os.Getenv("MQ_URL")
	// log.Print("MQ_URL: ", MQ_URL)

	if MQ_URL == "" {
		log.Print("⛔ MQ configs not found")
	}

	if client, err := nats.Connect(MQ_URL); err != nil {
		log.Print("⛔ MQ ", err)
	} else {
		MQClient = client
		log.Print("👍 MQ connected")
	}
}
