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

	//下面都是中间件
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOrigins:     []string{"http://localhost:5173"}, //允许前端哪个端口来访问我
		ExposeHeaders:    []string{"X-Request-ID"},
	}))

	authapi := api.NewAuthApi()
	favouriteapi := api.NewFavouriteApi()
	vndbapi := api.NewVNDBApi()
	public := router.Group("/api")
	{
		public.POST("/auth/register", authapi.Register)
		public.POST("/auth/login", authapi.Login)
	}

	//详情看Gin的路由分组
	private := router.Group("/api")
	private.Use(middleware.JWTAuth())
	{
		private.GET("/user/profile", authapi.Profile)
		//只给搜索接口挂限流
		private.GET("/v1/vndb/search", middleware.Ratelimit(10, time.Minute), vndbapi.Search)

		//这个是URL吗，怎么看着这么怪？
		private.POST("/v1/favourite", favouriteapi.Create)
		private.GET("/v1/favourite", favouriteapi.List)
		private.DELETE("/v1/favourite/:vndb_id", favouriteapi.Delete)
	}
	router.Run("127.0.0.1:8080")
}
