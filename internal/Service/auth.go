package Service

import (
	"Gal-Finder/internal/global"
	"Gal-Finder/internal/middleware"
	"Gal-Finder/internal/model"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	Invailduserformat = errors.New("the username format error")
	ErrUserExists     = errors.New("the username already exists")
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResult struct {
	Token    string   `json:"token"`
	UserInfo UserInfo `json:"userInfo"`
}

type UserInfo struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	UserID   uint   `json:"userID"`
}

func (a *AuthService) Login(req LoginRequest) (loginResult LoginResult, err error) {
	var user model.User
	if err := global.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return LoginResult{}, Invailduserformat
		}
		return LoginResult{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return LoginResult{}, Invailduserformat
	}

	//在登录的时候生成一个token,一般一个token带有签发时间，过期时间，签发者
	claims := middleware.Claims{
		UserID: user.ID, //一定要到这才用吗?
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	//用密钥把claims签名，变成token字符串
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtSecret)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{Token: tokenString, UserInfo: UserInfo{
		UserID:   user.ID,
		Nickname: user.Nickname,
		Username: user.Username,
	}}, nil
}
