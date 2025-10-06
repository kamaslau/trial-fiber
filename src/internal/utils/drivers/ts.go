package drivers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
		log.Print("⚠️ TS configs not found")
		return
	}

	// Confs
	configs := influxdb3.ClientConfig{
		Host:     host,
		Token:    token,
		Database: database,
	}

	// TODO debug this, won't break even fake config env is provided
	if client, err := influxdb3.New(configs); err != nil {
		log.Printf("🛑 TS connection failed: %v", err)
		return
	} else {
		TSClient = client
		log.Print("✅ TS connected")
	}

	logOnTSConnected()
}

// logOnTSConnected 写日志：时序存储连接成功
func logOnTSConnected() {
	point := influxdb3.NewPoint(
		"connection_log",
		map[string]string{"type": "ts"},
		map[string]interface{}{
			"status":  "connected",
			"message": fmt.Sprintf("Time series storage connection succeed at %s", time.Now().Format(time.RFC3339)),
		},
		time.Now(),
	)

	if err := TSClient.WritePoints(context.Background(), []*influxdb3.Point{point}); err != nil {
		log.Printf("└ Failed to write point: %v", err)
		return
	}
	log.Printf("└ Succeed to write point")
}
