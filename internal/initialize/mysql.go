package initialize

import (
	"Gal-Finder/internal/global"
	"Gal-Finder/internal/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() {
	dsn := "root:123456@tcp(127.0.0.1:3307)/galfinder?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect mysql fail %v", err)

	}

	if err := db.AutoMigrate(&model.User{}, &model.Favourite{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}
	global.DB = db
}
