package drivers

import (
	"log"
	"os"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
)

// TSClient 连接单例
var TSClient *influxdb3.Client

// ConnectTS 建立连接
func ConnectTS() {
	host := os.Getenv("TS_HOST")
	token := os.Getenv("TS_TOKEN")
	database := os.Getenv("TS_DATABASE")

	if host == "" || token == "" || database == "" {
		log.Print("⛔ TS configs not found")
		return
	}

	// Confs
	configs := influxdb3.ClientConfig{
		Host:     host,
		Token:    token,
		Database: database,
	}
	if client, err := influxdb3.New(configs); err != nil {
		log.Print("⛔ TS ", err)
	} else {
		TSClient = client
		log.Print("👍 TS connected")
	}
}
