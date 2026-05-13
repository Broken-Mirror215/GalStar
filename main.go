package main

import (
	"Gal-Finder/internal/api"
	"time"

	"Gal-Finder/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOrigins:     []string{"http://localhost:8080"},
		ExposeHeaders:    []string{"X-Request-ID"},
	}))
	authapi := api.NewAuthApi()
	favouriteapi := api.NewFavouriteApi()
	vndbapi := api.NewVNDBApi()
	public := router.Group("/api")
	{
		public.POST("/auth/Register", authapi.Register)
		public.POST("/auth/Login", authapi.Login)
	}

	//详情看Gin的路由分组
	private := router.Group("/api")
	private.Use(middleware.JWTAuth())
	{
		private.GET("/user/profile", authapi.Profile)
		private.GET("/v1/vndb/search", vndbapi.Search)

		//这个是URL吗，怎么看着这么怪？
		private.POST("/v1/favourite", favouriteapi.Create)
		private.GET("/v1/favourite", favouriteapi.List)
		private.DELETE("/v1/favourite/:vndb_id", favouriteapi.Delete)
	}
	router.Run("127.0.0.1:8080")
}
