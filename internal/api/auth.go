package api

import (
	"Gal-Finder/internal/global"
	"Gal-Finder/internal/middleware"
	"Gal-Finder/internal/model"
	"Gal-Finder/internal/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	var user model.User
	if err := global.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		response.Fail(c, 401, 401, "invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Fail(c, 401, 401, "invalid username or password")
		return
	}

	UserID := user.ID
	//在登录的时候生成一个token
	//一般一个token带有签发时间，过期时间，签发者
	//这就是定义中间件
	claims := middleware.Claims{
		UserID: UserID,
		RegisteredClaims: jwt.RegisteredClaims{
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
