package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"fmt"
)

var secretKey = []byte("secret-key")

func CreateToken(name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": name,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "error", err
	}
	fmt.Println(tokenString)
	return tokenString, nil
}
