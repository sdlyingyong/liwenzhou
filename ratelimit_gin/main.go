package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ratelimit2 "github.com/juju/ratelimit"
	ratelimit1 "go.uber.org/ratelimit"
)

func main() {
	r := gin.Default()
	r.GET("/ping", rateLimit1(), pingHandler)
	r.GET("/hei", rateLimit2(), heiHandler)
	r.Run(":8080")
}

//基于漏桶的限流中间件
//一个水桶按照固定的速率向下滴水,也就是固定速度执行
func rateLimit1() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 每秒放行的次数
		rl := ratelimit1.New(1)
		//下个雨滴的时间小于当前时间,需要等待
		next := rl.Take()
		if next.Sub(time.Now()) > 0 {
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		//放行
		c.Next()
	}
}

//基于令牌桶的限流中间件
func rateLimit2() func(c *gin.Context) {
	//设置令牌生成速度
	bucket := ratelimit2.NewBucket(time.Second*2, 1)
	return func(c *gin.Context) {
		//获取令牌发现用完,限制
		//Take()  拿到令牌,没有可以欠账
		//TakeAvailable() 只有现存令牌才能取出,没有令牌会等待
		if bucket.TakeAvailable(1) == 1 {
			//放行
			c.Next()
			return
		}
		c.String(http.StatusOK, "rate limit...")
		//终止请求
		c.Abort()
		return
	}
}

func pingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
	return
}

func heiHandler(c *gin.Context) {
	c.String(http.StatusOK, "hei")
	return
}
