package hooks

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"net/http"
	"strings"
	"fmt"
);

type MyCustomClaims struct {
	UserID uint `json:"userId"`
	jwt.StandardClaims
}

func Descrypt(Authtoken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpectedsigning method: %v", token.Header["alg"])
			}
			return []byte(Authtoken), nil
		})

		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			// Store user information in the context
			c.Set("userID", claims.UserID)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		
	}
}