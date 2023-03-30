package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/plant_disease_detection/internal/auth"
)

func Protect(c *gin.Context) {
	auth_header := c.GetHeader("Authorization")
	if len(auth_header) == 0 {
		log.Println("no token given")
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
		return
	}

	tokenString := strings.Split(auth_header, "Bearer")[1][1:]

	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
		c.Set("user", claims)
		c.Set("authtoken", token.Raw)
	} else {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
	}
}
