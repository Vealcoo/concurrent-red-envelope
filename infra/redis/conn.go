package redis

import (
	"context"
	"log"
	"test/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func init() {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Setting.GetString("redis.addr"),
		Username: "",
		Password: "",
		DB:       0,
		// PoolSize:     200,             // 最大連線數，根據伺服器資源調整
		// MinIdleConns: 20,              // 最小閒置連線數
		// PoolTimeout:  time.Second * 5, // 等待連線的超時時間
		// MaxRetries:   3,               // 失敗重試次數
		// DialTimeout:  time.Second * 5, // 建立連線超時
		// ReadTimeout:  time.Second * 3, // 讀取超時
		// WriteTimeout: time.Second * 3, // 寫入超時
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatal(err)
	}

	RedisClient = client
}
