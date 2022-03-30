package main

import (
	"middleware/controller"
	"middleware/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)
	router.GET("/secure-welcome", middlewares.AuthMiddleware(), controller.Welcome)
	router.GET("/unsecure-welcome", controller.Welcome)
	router.Run(":8091")
}
