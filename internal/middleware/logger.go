package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		//这里直接next的原因是要等后面的弄完，所以logger是一种前后包裹的中间件
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		ClientIp := c.ClientIP()

		Requestid, _ := c.Get("requestID")
		user, _ := c.Get("user")

		fmt.Printf("status=%d\n method=%s\n path=%s\n ip=%s\n requestid=%v\n user=%v \nlatency=%s", status, method, path, ClientIp, Requestid, user, latency)
	}

}
