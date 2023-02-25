package formtoken

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"
)

func Test_FormToken(t *testing.T) {
	var (
		token = GetFormToken()
		a     = 10
		wg    = sync.WaitGroup{}
	)
	wg.Add(a)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	for i := 0; i < a; i++ {
		go func(i int) {
			defer wg.Done()
			if Access(ctx, token) {
				log.Printf("进程%d执行成功", i)
			} else {
				cancel()
				log.Printf("进程%d执行失败", i)
			}
		}(i)
	}
	wg.Wait()
}
