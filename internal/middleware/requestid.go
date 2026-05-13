package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		//_ = requestID
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}

		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
