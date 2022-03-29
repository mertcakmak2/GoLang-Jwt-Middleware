package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	router := gin.Default()
	router.POST("/register", register)
	router.POST("/login", login)
	router.GET("/secure-welcome", authMiddleware(), welcome)
	router.GET("/unsecure-welcome", welcome)
	router.Run(":8091")
}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, "welcome")
}

func login(c *gin.Context) {

	fmt.Println(users)

	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userHashPassword, isPresent := users[creds.Username]
	if !isPresent {
		c.JSON(http.StatusUnauthorized, "Wrong username or password")
		return
	}

	if match := CheckPasswordHash(creds.Password, userHashPassword); !match {
		c.JSON(http.StatusUnauthorized, "Wrong username or password")
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error")
		return
	}

	c.JSON(http.StatusOK, tokenString)
}

func register(c *gin.Context) {

	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, _ := HashPassword(creds.Password)
	users[creds.Username] = hash

	creds.Password = hash
	c.JSON(http.StatusCreated, creds)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" || !strings.Contains(authorization, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
			c.Abort()
			return
		}

		token := strings.Split(authorization, "Bearer ")[1]
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, "Authentication failed")
				c.Abort()
				return
			}
		}
		if !tkn.Valid {
			c.JSON(http.StatusUnauthorized, "Authentication failed")
			c.Abort()
			return
		}
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var users = map[string]string{
	// "user1": "password1",
	// "user2": "password2",
}

var jwtKey = []byte("my_secret_key")
