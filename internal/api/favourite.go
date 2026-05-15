package api

import (
	"Gal-Finder/internal/response"
	"Gal-Finder/internal/service"
	"errors"

	"github.com/gin-gonic/gin"
)

type FavouriteApi struct {
	FavorService *service.FavoriteService
}

func NewFavouriteApi() *FavouriteApi {
	return &FavouriteApi{
		FavorService: service.NewFavoriteService(),
	}
}

// 自己定义的后端接口喵,尽量和数据库对齐

func (a *FavouriteApi) Create(c *gin.Context) {
	userIDValue, _ := c.Get("userID")
	userID := userIDValue.(uint)
	var req service.CreateFavouriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 400, "invaild request")
		return
	}

	err := a.FavorService.Create(userID, req)
	if errors.Is(err, service.ErrfavorNotExisted) {
		response.Fail(c, 400, 400, "你已经收藏了喵")
		return
	}
	if err != nil {
		response.Fail(c, 500, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"userId": userID,
	})
}

func (a *FavouriteApi) List(c *gin.Context) {
	userIDValue, _ := c.Get("userID")
	userID := userIDValue.(uint)
	favourites, err := a.FavorService.List(userID)
	if err != nil {
		response.Fail(c, 500, 500, "query favourites error")
		return
	}

	response.Success(c, gin.H{
		"list":   favourites,
		"userID": userID,
	})
}

func (a *FavouriteApi) Delete(c *gin.Context) {
	userIDValue, _ := c.Get("userID")
	userID := userIDValue.(uint)
	vndbID := c.Param("vndb_id")

	err := a.FavorService.Delete(userID, vndbID)
	if errors.Is(err, service.ErrorfavorNotFound) {
		response.Fail(c, 404, 404, "favorite not found")
		return
	}
	if err != nil {
		response.Fail(c, 500, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"userId": userID,
	})
}
