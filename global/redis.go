package global

import (
	"iris-project/config"

	"github.com/go-redis/redis/v7"
)

// Redis 连接
var Redis *redis.Client

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
		PoolSize: config.Redis.PoolSize,
	})

	_, err := Redis.Ping().Result()
	if err != nil {
		panic(err)
	}

}
