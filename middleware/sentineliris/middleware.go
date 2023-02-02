package sentineliris

import (
	"net/http"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/kataras/iris/v12/context"
)

// SentinelMiddleware returns new context.Handler
// Default resource name is {method}:{path}, such as "GET:/api/users/:id"
// Default block fallback is returning 429 code
// Define your own behavior by setting options
func SentinelMiddleware(opts ...Option) context.Handler {
	options := evaluateOptions(opts)
	return func(c *context.Context) {
		// resourceName := c.Request.Method + ":" + c.FullPath()
		resourceName := c.Method() + ":" + c.GetCurrentRoute().ResolvePath()
		// log.Println("资源名称", resourceName)

		if options.resourceExtract != nil {
			resourceName = options.resourceExtract(c)
		}

		entry, err := sentinel.Entry(
			resourceName,
			sentinel.WithResourceType(base.ResTypeWeb),
			sentinel.WithTrafficType(base.Inbound),
		)

		if err != nil {
			if options.blockFallback != nil {
				options.blockFallback(c)
			} else {
				c.StopWithStatus(http.StatusTooManyRequests)
			}
			return
		}

		defer entry.Exit()
		c.Next()
	}
}
