package middleware

import (
	"Gal-Finder/internal/response"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	Count     int
	ResetTime time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex //这看起来像一个异步锁
)

func Ratelimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		if UserID, ok := c.Get("userID"); ok {
			key = fmt.Sprintf("user %v", UserID)
		}

		now := time.Now()

		mu.Lock()
		v, ok := visitors[key]
		if !ok || now.After(v.ResetTime) {
			visitors[key] = &visitor{
				Count:     1,
				ResetTime: now.Add(window), //这是什么？？
			}
			mu.Unlock()
			c.Next()
			return
		}

		if v.Count >= limit {
			mu.Unlock()
			response.Fail(c, 429, 429, "too many requests")
			c.Abort()
			return
		}
		v.Count++
		mu.Unlock()
		c.Next()
	}

}
