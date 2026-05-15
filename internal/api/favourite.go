package api

import (
	"Gal-Finder/internal/global"
	"Gal-Finder/internal/model"
	"Gal-Finder/internal/response"

	"github.com/gin-gonic/gin"
)

type FavouriteApi struct{}

func NewFavouriteApi() *FavouriteApi {
	return &FavouriteApi{}
}

// 自己定义的后端接口喵,尽量和数据库对齐
type CreateFavouriteRequest struct {
	VNDBID       string  `json:"VNDBID" binding:"required"`
	Title        string  `json:"Title" binding:"required"`
	ImageUrl     string  `json:"ImageUrl"`
	ThumbnailUrl string  `json:"ThumbnailUrl"`
	Rating       float64 `json:"Rating"`
	Released     string  `json:"Released"`
}

func (a *FavouriteApi) Create(c *gin.Context) {
	userIDValue, _ := c.Get("userID")
	userID := userIDValue.(uint)

	var req CreateFavouriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 400, "invaild request")
		return
	}

	favorite := model.Favourite{
		VNDBID:       req.VNDBID,
		UserID:       userID,
		Title:        req.Title,
		ImageUrl:     req.ImageUrl,
		Rating:       req.Rating,
		Released:     req.Released,
		ThumbnailUrl: req.ThumbnailUrl,
	}

	if err := global.DB.Create(&favorite).Error; err != nil {
		response.Fail(c, 400, 400, "你已经收藏了喵")
		return
	}

	response.Success(c, gin.H{
		"userID":    userID,
		"favourite": 1,
	})
}

func (a *FavouriteApi) List(c *gin.Context) {
	userIDValue, _ := c.Get("userID")
	userID := userIDValue.(uint)

	var favourites []model.Favourite
	if err := global.DB.Where("user_id = ?", userID).Find(&favourites).Error; err != nil {
		response.Fail(c, 500, 500, "query List failed")
		return
	}

	response.Success(c, gin.H{
		"userID": userID,
		"list":   favourites,
	})
}

func (a *FavouriteApi) Delete(c *gin.Context) {
	userIDValue, _ := c.Get("userID")
	userID := userIDValue.(uint)
	vndbID := c.Param("vndb_id")

	result := global.DB.Where("user_id  = ? AND vndb_id = ?", userID, vndbID).Delete(&model.Favourite{})
	if result.Error != nil {
		response.Fail(c, 500, 500, "Delete failed")
		return
	}

	if result.RowsAffected == 0 {
		response.Fail(c, 404, 404, "favourite not found")
		return
	}

	response.Success(c, gin.H{
		"VNDBID": vndbID,
	})
}
