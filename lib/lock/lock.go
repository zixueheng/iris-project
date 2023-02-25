package lock

/*
import (
	"iris-project/lib/cache"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v7"
)

const (
	letters     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lockCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	// lockCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
	// 	redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
	// 	return "OK"
	// else
	// 	return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
	// end`

	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`

	// delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
	// return redis.call("DEL", KEYS[1])
	// else
	// return 0
	// end`

	randomLen = 16
	// 默认超时时间，防止死锁
	tolerance       = 500 // milliseconds
	millisPerSecond = 1000
)

var (
	redisDb *redis.Client
)

// A RedisLock is a redis lock.
type RedisLock struct {
	// redis客户端
	// store *redis.Client
	// 超时时间
	seconds uint32
	// 锁key
	key string
	// 锁value，防止锁被别人获取到
	id string
}

func init() {
	rand.Seed(time.Now().UnixNano())
	redisDb = cache.GetCacheInstance()
}

// NewRedisLock returns a RedisLock.
func NewRedisLock(key string) *RedisLock {
	return &RedisLock{
		key: key,
		// 获取锁时，锁的值通过随机字符串生成
		// 实际上go-zero提供更加高效的随机字符串生成方式
		// 见core/stringx/random.go：Randn
		id: randomStr(randomLen),
	}
}

// Acquire acquires the lock.
// 加锁
func (rl *RedisLock) Acquire() (bool, error) {
	// 获取过期时间
	seconds := atomic.LoadUint32(&rl.seconds)
	// 默认锁过期时间为500ms，防止死锁
	resp, err := redisDb.Eval(lockCommand, []string{rl.key}, []string{
		rl.id, strconv.Itoa(int(seconds)*millisPerSecond + tolerance),
	})
	t := strconv.Itoa(rl.timeoutMs)
	redisDb.Eval(lockCommand, []string{rl.key}, []string{rl.id, t}).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		// logx.Errorf("Error on acquiring lock for %s, %s", rl.key, err.Error())
		return false, err
	} else if resp == nil {
		return false, nil
	}

	reply, ok := resp.(string)
	if ok && reply == "OK" {
		return true, nil
	}

	// logx.Errorf("Unknown reply when acquiring lock for %s: %v", rl.key, resp)
	return false, nil
}

// Release releases the lock.
// 释放锁
func (rl *RedisLock) Release() (bool, error) {
	resp, err := redisDb.Eval(delCommand, []string{rl.key}, []string{rl.id})
	if err != nil {
		return false, err
	}

	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}

	return reply == 1, nil
}

// SetExpire sets the expire.
// 需要注意的是需要在Acquire()之前调用
// 不然默认为500ms自动释放
func (rl *RedisLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.seconds, uint32(seconds))
}

func randomStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
*/
