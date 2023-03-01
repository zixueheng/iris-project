/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-02-02 10:57:08
 * @LastEditTime: 2023-02-03 10:24:28
 */
package middleware

import (
	"errors"
	"iris-project/app"
	"iris-project/middleware/sentineliris"
	"log"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// Sentinel 限流
func Sentinel() context.Handler {
	err := sentinel.InitWithConfigFile("config/sentinel.yml")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := system.LoadRules([]*system.Rule{
		{
			MetricType:   system.InboundQPS,
			TriggerCount: 500, // qps 500
			Strategy:     system.BBR,
		},
	}); err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}
	return sentineliris.SentinelMiddleware(
		sentineliris.WithBlockFallback(func(ctx *context.Context) {
			app.ResponseProblemHTTPCode(ctx, iris.StatusTooManyRequests, errors.New("太多请求，请稍后重试"))
		}),
	)
}
