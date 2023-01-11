package utils

import (
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func GetUserId(tokenString string) string {
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, _ := jwt.Parse(tokenString, nil)
	claims, _ := token.Claims.(jwt.MapClaims)
	return claims["sub"].(string)
}
