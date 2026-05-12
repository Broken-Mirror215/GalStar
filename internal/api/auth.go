package api

import (
	"Gal-Finder/internal/response"

	"github.com/gin-gonic/gin"
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
}
func (a *AuthApi) Register(c *gin.Context) {
	var req Register
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 401, 401, err.Error())
		return
	}
}
