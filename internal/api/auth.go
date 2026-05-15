package api

import (
	"Gal-Finder/internal/response"
	"Gal-Finder/internal/service"
	"errors"

	"github.com/gin-gonic/gin"
)

type AuthApi struct {
	AuthService *service.AuthService
}

func NewAuthApi() *AuthApi {
	return &AuthApi{
		AuthService: service.NewAuthService(),
	}
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
	var req service.LoginRequest
	//把前端登录送来的数据放进去
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 401, 401, err.Error())
		return
	}

	result, err := a.AuthService.Login(req)
	if errors.Is(err, service.Invailduserformat) {
		response.Fail(c, 500, 500, "Invaild username or password")
		return
	}

	if err != nil {
		response.Fail(c, 500, 500, err.Error())
		return
	}
	response.Success(c, result)
}

func (a *AuthApi) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 400, "the user name has existed")
		return
	}
	user, err := a.AuthService.Register(req)
	if err != nil {
		response.Fail(c, 500, 500, err.Error())
		return
	}
	response.Success(c, user)
}

func (a *AuthApi) Profile(c *gin.Context) {
	userId, _ := c.Get("userID")
	response.Success(c, gin.H{
		"userID": userId,
	})
}
