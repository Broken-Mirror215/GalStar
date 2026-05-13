package api

import (
	"Gal-Finder/internal/middleware"
	"Gal-Finder/internal/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthApi struct{}

func NewAuthApi() *AuthApi {
	return &AuthApi{}
}

type Register struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *AuthApi) Login(c *gin.Context) {
	var req Login
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 401, 401, err.Error())
		return
	}

	//在登录的时候生成一个token
	//一般一个token带有签发时间，过期时间，签发者
	UserID := uint(1)
	claims := middleware.Claims{
		UserID: UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			//这是什么
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	//用密钥把claims签名，变成token字符串
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("change-this-secret"))
	if err != nil {
		response.Fail(c, 500, 500, err.Error())
		return
	}
	response.Success(c, gin.H{
		"token": tokenString,
		"user": gin.H{
			"userID": UserID,
		},
	})
}

func (a *AuthApi) Register(c *gin.Context) {
	var req Register
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 401, 401, err.Error())
		return
	}

	response.Success(c, gin.H{
		"username": req.Username,
		"nickname": req.Nickname,
	})
}

func (a *AuthApi) Profile(c *gin.Context) {
	userId, _ := c.Get("userID")
	response.Success(c, gin.H{
		"userID": userId,
	})
}
