package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.Run("127.0.0.1:8080")
}
