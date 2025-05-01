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
	CACHE_URL := os.Getenv("CACHE_URL")
	// log.Print("CACHE_URL: ", CACHE_URL)

	if CACHE_URL == "" {
		log.Print("⛔ Cache configs not found")
		return
	}

	if opts, err := redis.ParseURL(CACHE_URL); err != nil {
		log.Print("⛔ Cache ", err)
	} else {
		CacheClient = redis.NewClient(opts)
		log.Print("👍 Cache connected")
	}
}
