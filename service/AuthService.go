package service

import (
	"errors"
	"fmt"
	"middleware/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var users = map[string]string{}
var jwtKey = []byte("my_secret_key")

func Login(creds model.Credentials) (string, error) {
	userHashPassword, isPresent := users[creds.Username]
	if !isPresent {
		return "", errors.New("wrong username or password")
	}

	if match := checkPasswordHash(creds.Password, userHashPassword); !match {
		return "", errors.New("wrong username or password")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &model.Claims{
		Username: creds.Username,
		Roles:    []string{"ADMIN", "SUPERADMIN"},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", errors.New("error")
	}
	return tokenString, nil
}

func Register(creds *model.Credentials) (err error) {

	hash, _ := hashPassword(creds.Password)
	users[creds.Username] = hash

	creds.Password = hash
	fmt.Println(users)
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
