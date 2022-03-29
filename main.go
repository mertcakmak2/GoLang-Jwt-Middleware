package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/middleware", ensureLoggedIn(), welcome)
	router.GET("/without-middleware", welcome)

	router.Run(":8091")
}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, "welcome")
}

func ensureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		jwt := strings.Split(authorization, "Bearer ")[1]

		fmt.Println("ensure logged in middleware jwt :", jwt)

		if jwt != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
			c.Abort()
			return
		}
	}
}
