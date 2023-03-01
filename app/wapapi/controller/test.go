/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-03-19 14:01:28
 * @LastEditTime: 2023-02-27 14:40:21
 */
package controller

import (
	"iris-project/app"
	"iris-project/lib/lock"
	"log"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Test 控制器
type Test struct {
	Ctx iris.Context
}

// BeforeActivation 前置方法
func (c *Test) BeforeActivation(b mvc.BeforeActivation) {
}

// GetTest
func (c *Test) GetTest() {
	c.Ctx.Application().Logger().Info("测试LOG")
	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// GetTestLock
func (c *Test) GetTestLock() {
	c.Ctx.Application().Logger().Info("测试LOG")
	var lock = lock.NewRedisLockWithCtx(c.Ctx, "lock-test")
	if err := lock.Lock(); err != nil {
		c.Ctx.Application().Logger().Warn(time.Now().UnixMilli(), "失败")
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}
	c.Ctx.Application().Logger().Info(time.Now().UnixMilli(), "成功")
	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

func (c *Test) GetTestLock2() {
	var (
		a  = 10
		wg = sync.WaitGroup{}
	)
	wg.Add(a)

	for i := 0; i < a; i++ {
		go func(i int) {
			defer wg.Done()
			var lock = lock.NewRedisLockWithCtx(c.Ctx.Clone(), "lock-test")
			if err := lock.Lock(); err != nil {
				log.Printf("%d进程%d执行失败", time.Now().UnixMilli(), i)
			} else {
				log.Printf("%d进程%d执行成功", time.Now().UnixMilli(), i)
			}
		}(i)
	}
	wg.Wait()
}
