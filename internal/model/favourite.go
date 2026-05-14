package model

import "gorm.io/gorm"

type Favourite struct {
	gorm.Model
	UserID       uint   `gorm:"not null;index;uniqueIndex:idx_user_vndb"`
	VNDBID       string `gorm:"type:varchar(32);not null;index:idx_user_vndb"`
	Title        string `gorm:"type:varchar(255);not null"`
	ImageUrl     string `gorm:"not null;type:text"`
	ThumbnailUrl string `gorm:"type:text"`
	Rating       int
	Released     string `gorm:"type:varchar(32);not null"`
}
