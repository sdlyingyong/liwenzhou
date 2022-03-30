package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

//限制请求速度
//每个请求从桶里获取令牌,没有就直接返回,有才放行
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	//gin 中间件返回gin context的方法
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		//如果取不到令牌就返回响应
		if bucket.TakeAvailable(1) == 0 {
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		//获取到令牌就放行
		c.Next()
	}
}
