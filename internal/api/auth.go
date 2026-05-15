package api

import (
	"Gal-Finder/internal/global"
	"Gal-Finder/internal/model"
	"Gal-Finder/internal/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	//把前端登录送来的数据放进去
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 401, 401, err.Error())
		return
	}

}

func (a *AuthApi) Register(c *gin.Context) {
	var req Register
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 400, err.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, 500, 500, "hash password failed")
		return
	}

	user := model.User{
		Username:     req.Username,
		Nickname:     req.Nickname,
		PasswordHash: string(hash),
	}

	if err := global.DB.Create(&user).Error; err != nil {
		response.Fail(c, 500, 500, "the user name already exists")
		return
	}

	response.Success(c, gin.H{
		"userID":   user.ID, //哪里出现的？
		"username": user.Username,
		"nickname": user.Nickname,
	})
}

func (a *AuthApi) Profile(c *gin.Context) {
	userId, _ := c.Get("userID")
	response.Success(c, gin.H{
		"userID": userId,
	})
}
