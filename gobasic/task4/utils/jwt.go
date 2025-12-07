package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecretKey := "1qaz@WSX"
	return token.SignedString([]byte(jwtSecretKey))
}

func ParseToken(tokenString string) (*jwt.MapClaims, error) {
	jwtSecretKey := "1qaz@WSX"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}
	return nil, fmt.Errorf("无效的token")
}
