package middleware

import (
	"Gal-Finder/internal/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("change-this-secret")

type Claims struct {
	UserID uint `json:"userid"`
	jwt.RegisteredClaims
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(c, 401, 401, "missing token")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Fail(c, 400, 401, "token format error")
		}
		token, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			response.Fail(c, 400, 401, "invalid token")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(*Claims)
		if !ok {
			response.Fail(c, 400, 401, "invalid token")
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
