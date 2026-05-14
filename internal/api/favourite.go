package api

import (
	"Gal-Finder/internal/response"

	"github.com/gin-gonic/gin"
)

type FavouriteApi struct{}

func NewFavouriteApi() *FavouriteApi {
	return &FavouriteApi{}
}

// 自己定义的后端接口喵
type CreateFavouriteRequest struct {
	VNDBID       string `json:"VNDBID" binding:"required"`
	Title        string `json:"Title" binding:"required"`
	ImageUrl     string `json:"ImageUrl"`
	ThumbnailUrl string `json:"ThumbnailUrl"`
	Rating       int    `json:"Rating"`
	Released     string `json:"Released"`
}

func (a *FavouriteApi) Create(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req CreateFavouriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 400, "invaild request")
		return
	}
	response.Success(c, gin.H{
		"userID":    userID,
		"favourite": 1,
	})
}

func (a *FavouriteApi) List(c *gin.Context) {
	userID, _ := c.Get("userID")
	response.Success(c, gin.H{
		"userID": userID,
		"list":   []interface{}{}, //返回一个空的接口类型
	})
}

func (a *FavouriteApi) Delete(c *gin.Context) {
	userID, _ := c.Get("userID")
	vndbID := c.Param("vndb_id")

	response.Success(c, gin.H{
		"userID": userID,
		"VNDBID": vndbID,
	})
}
