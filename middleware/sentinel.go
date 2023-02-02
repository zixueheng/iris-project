package middleware

import (
	"errors"
	"iris-project/app"
	"iris-project/middleware/sentineliris"
	"log"

	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// Sentinel 限流
func Sentinel() context.Handler {
	// err := sentinel.InitDefault()
	// if err != nil {
	// 	log.Fatal(err)
	// }
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
