/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-23 10:16:29
 * @LastEditTime: 2021-04-28 14:46:46
 */
package cache

import (
	"iris-project/global"
	"time"

	"github.com/go-redis/redis/v7"
)

// GetCacheInstance 返回缓存对象
func GetCacheInstance() *redis.Client {
	return global.Redis
}

// Get 按键取缓存值
func Get(key string) (string, error) {
	return GetCacheInstance().Get(key).Result()
}

// Set 按键保存缓存值
func Set(key string, value interface{}, expiration time.Duration) {
	GetCacheInstance().Set(key, value, expiration)
}

// Del 按多个键删除缓存值
func Del(keys ...string) {
	GetCacheInstance().Del(keys...)
}
