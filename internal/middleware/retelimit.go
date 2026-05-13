package middleware

import (
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

func Ratelimie(limit int, window time.Duration) gin.HandlerFunc {

}
