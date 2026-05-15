package model

import "gorm.io/gorm"

type Favourite struct {
	gorm.Model
	UserID       uint    `gorm:"not null;index;uniqueIndex:idx_user_vndb"`
	VNDBID       string  `gorm:"type:varchar(32);not null;uniqueIndex:idx_user_vndb"` //保证不重复收藏！
	Title        string  `gorm:"type:varchar(255);not null"`
	ImageUrl     string  `gorm:"not null;type:text"`
	ThumbnailUrl string  `gorm:"type:text"`
	Rating       float64 `gorm:"type:float"`
	Released     string  `gorm:"type:varchar(32);not null"`
}
