package hooks

import (
	"time"
	"os"
   
	"github.com/golang-jwt/jwt/v5"
   )

func CreateToken(user any) (any, error) {
	secret := os.Getenv("JWT_SECRET");

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user": user,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
	
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}