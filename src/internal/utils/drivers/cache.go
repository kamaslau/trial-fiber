package drivers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9" // https://redis.uptrace.dev/guide/
)

// CacheClient 连接单例
var CacheClient *redis.Client

// ConnectCache 建立连接
func ConnectCache() {
	CACHE_URL := os.Getenv("CACHE_URL")
	// log.Print("CACHE_URL: ", CACHE_URL)

	if CACHE_URL == "" {
		log.Print("⚠️ Cache configs not found")
		return
	}

	if opts, err := redis.ParseURL(CACHE_URL); err != nil {
		log.Print("🛑 Cache ", err)
	} else {
		CacheClient = redis.NewClient(opts)
		log.Print("✅ Cache connected")
	}

	// Optional Ops
	logOnCacheConnected()
	// TODO Load from database
	// - revoked tokens, i.e. "token-revoke-user:"+id
}

// logOnCacheConnected 写日志：连接成功
func logOnCacheConnected() {
	if err := CacheSet("Cache connected", fmt.Sprintf("This is an auto generated message on cache connection succeed at %s", time.Now().Format(time.RFC3339)), 0); err != nil {
		log.Printf("└ Failed to insert data: %v", err)
		return
	}

	log.Printf("└ Succeed to insert data key: %s", "Cache connected")
}

func CacheSet(key string, value any, expiration time.Duration) error {
	ctx := context.Background()

	if err := CacheClient.Set(ctx, key, value, expiration).Err(); err != nil {
		return err
	}

	return nil
}

func CacheGet(key string) (string, error) {
	ctx := context.Background()

	value, err := CacheClient.Get(ctx, key).Result()
	if err == redis.Nil {
		// log.Printf("key %s does not exist", key)
	} else if err != nil {
		// log.Printf("CacheGet error: %v", err)
	} else {
		// log.Printf("Got cache: for key %s, value is %v", key, value)
	}

	return value, err
}
