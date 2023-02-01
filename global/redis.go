/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:27:34
 * @LastEditTime: 2022-05-05 16:14:15
 */
package global

import (
	"iris-project/app/config"
	"log"
	"time"

	"github.com/go-redis/redis/v7"
)

// Redis 连接
var Redis *redis.Client

func init() {
	opt := &redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
		PoolSize: config.Redis.PoolSize,
	}
	Redis = redis.NewClient(opt)

	_, err := Redis.Ping().Result()
	if err != nil {
		for {
			time.Sleep(60 * time.Second) // 等待60秒再次连接
			log.Println("再次连接Redis")
			Redis = redis.NewClient(opt)
			_, err := Redis.Ping().Result()
			if err == nil {
				break
			} else {
				log.Println(err.Error())
			}
		}
		// panic(err)
	}

}
