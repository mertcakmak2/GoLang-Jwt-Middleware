package middlewares

import (
	"middleware/model"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("my_secret_key")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" || !strings.Contains(authorization, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
			c.Abort()
			return
		}

		token := strings.Split(authorization, "Bearer ")[1]
		claims := &model.Claims{}

		jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"Message": err.Error()})
			c.Abort()
			return
		}
		if !jwtToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
			c.Abort()
			return
		}
	}
}
