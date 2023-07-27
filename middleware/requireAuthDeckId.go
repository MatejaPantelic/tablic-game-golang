package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"main.go/constants"
)

func CheckAuthTokenDeckId(c *gin.Context) {
	//extract parameter from request
	deckId := c.Param("deckId")

	//get authorization header
	authHeader := c.GetHeader("Authorization")

	//check if auth header exist
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.AUTH_HEADER_MISSING})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// split parts from auth header to get payload
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.INVALID_TOKEN})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// get payload from splited parts
	tokenString := authHeaderParts[1]

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["Alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	// validation if player is in game, validation if token has expired
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if deckId != claims["deck_id"] {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.FORBIDDEN_ACCESS})
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.TOKEN_EXPIRED})
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.INVALID_TOKEN})
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
