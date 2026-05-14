package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Nickname     string `gorm:"type:varchar(255)"`
}
