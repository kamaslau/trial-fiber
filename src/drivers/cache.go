package drivers

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// CacheClient 连接单例
var CacheClient *redis.Client

// ConnectCache 建立连接
func ConnectCache() {
	REDIS_URL := os.Getenv("REDIS_URL")
	// log.Print("REDIS_URL: ", REDIS_URL)

	opts, err := redis.ParseURL(REDIS_URL)
	if err != nil {
		panic(err)
	} else {
		CacheClient = redis.NewClient(opts)
		log.Print("Cache connected")
	}
}
