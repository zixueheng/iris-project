package lock

import (
	"context"
	"errors"
	"iris-project/lib/cache"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
)

const (

	// 先get获取，如果有就刷新ttl，没有再set。 这种是可重入锁，防止在同一线程中多次获取锁而导致死锁发生。
	lockCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
	redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
	return "OK"
else
	return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`

	// 删除。必须先匹配id值，防止A超时后，B马上获取到锁，A的解锁把B的锁删了
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end`

	letters   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomLen = 10
)

var (
	redisDb *redis.Client
	// 默认超时时间
	defaultTimeout = 500 * time.Millisecond
	// 重试间隔
	retryInterval = 10 * time.Millisecond
	// 上下文取消
	ErrContextCancel = errors.New("context cancel")
)

func init() {
	rand.Seed(time.Now().UnixNano())
	redisDb = cache.GetCacheInstance()
}

type RedisLock struct {
	ctx       context.Context
	timeoutMs int
	key       string
	id        string
}

func NewRedisLockWithCtx(ctx context.Context, key string) *RedisLock {
	timeout := defaultTimeout
	if deadline, ok := ctx.Deadline(); ok {
		timeout = deadline.Sub(time.Now())
	}

	rl := &RedisLock{
		ctx:       ctx,
		timeoutMs: int(timeout.Milliseconds()),
		key:       key,
		id:        randomStr(randomLen),
	}

	return rl
}

func NewRedisLock( /*ctx context.Context,*/ key string) *RedisLock {
	timeout := defaultTimeout
	// if deadline, ok := ctx.Deadline(); ok {
	// 	timeout = deadline.Sub(time.Now())
	// }
	rl := &RedisLock{
		// ctx:       ctx,
		timeoutMs: int(timeout.Milliseconds()),
		key:       key,
		id:        randomStr(randomLen),
	}

	return rl
}

func (rl *RedisLock) TryLock() (bool, error) {
	t := strconv.Itoa(rl.timeoutMs)
	resp, err := redisDb.Eval(lockCommand, []string{rl.key}, []string{rl.id, t}).Result()
	if err != nil || resp == nil {
		return false, nil
	}

	reply, ok := resp.(string)
	return ok && reply == "OK", nil
}

func (rl *RedisLock) Lock() error {
	for {
		select {
		case <-rl.ctx.Done():
			// log.Println("取消", time.Now().UnixMilli())
			return ErrContextCancel
		default:
			b, err := rl.TryLock()
			if err != nil {
				return err
			}
			if b {
				// log.Println("成功", time.Now().UnixMilli())
				return nil
			}
			// log.Println("重试", time.Now().UnixMilli())
			time.Sleep(retryInterval)
		}
	}
}

func (rl *RedisLock) Unlock() {
	redisDb.Eval(delCommand, []string{rl.key}, []string{rl.id}).Result()
}

func randomStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
