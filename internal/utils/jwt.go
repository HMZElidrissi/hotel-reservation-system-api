package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

var jwtSecret = []byte(viper.GetString("JWT_SECRET"))

func GenerateJWT(email, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.Claims(jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   expirationTime.Unix(),
	})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Claims, error) {
	var claims jwt.Claims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return &claims, nil
}
