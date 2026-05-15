package service

import (
	"Gal-Finder/internal/global"
	"Gal-Finder/internal/model"
	"errors"
)

var (
	ErrfavorNotExisted = errors.New("favor not existed")
	ErrorfavorNotFound = errors.New("favor not Found")
)

type FavoriteService struct{}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{}
}

type CreateFavouriteRequest struct {
	VNDBID       string  `json:"VNDBID" binding:"required"`
	Title        string  `json:"Title" binding:"required"`
	ImageUrl     string  `json:"ImageUrl"`
	ThumbnailUrl string  `json:"ThumbnailUrl"`
	Rating       float64 `json:"Rating"`
	Released     string  `json:"Released"`
}

func (a *FavoriteService) Create(userID uint, req CreateFavouriteRequest) error {
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
		return ErrfavorNotExisted
	}

	return nil
}

func (a *FavoriteService) List(userID uint) ([]model.Favourite, error) {
	var favourites []model.Favourite
	if err := global.DB.Where("user_id = ?", userID).Find(&favourites).Error; err != nil {
		return nil, err
	}

	return favourites, nil
}

func (a *FavoriteService) Delete(userID uint, vndbID string) error {
	result := global.DB.Where("user_id  = ? AND vndb_id = ?", userID, vndbID).Delete(&model.Favourite{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrorfavorNotFound
	}

	return nil
}
