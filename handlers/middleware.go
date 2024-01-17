// handlers/middleware.go
package handlers

import (
	"errors"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"

)

func AuthMiddleware() gin.HandlerFunc {
	var jwtSecret = []byte("your-secret-key")
	
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			ErrorResponse(c, http.StatusUnauthorized, errors.New("Authorization header is missing"))
			c.Abort()
			return
		}
		
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ErrorResponse(c, http.StatusUnauthorized, errors.New("Invalid header format or missing \"Bearer\" prefix"))
			c.Abort()
			return
		}
		
		tokenString := parts[1]
		
		if tokenString == "" {
			ErrorResponse(c, http.StatusUnauthorized, errors.New("Authorization header is missing"))
			c.Abort()
			return
		}
		
		// Verify and parse JWT token
		claims, err := ParseJWT(tokenString, jwtSecret)
		if err != nil {
			ErrorResponse(c, http.StatusUnauthorized, errors.New("Invalid or expired token"))
			c.Abort()
			return
		}
		
		// Set user ID from claims in context for later use
		c.Set("UserID", claims.UserID)
		
		c.Next()
	}
}

func ParseJWT(tokenString string, jwtSecret []byte) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	
	if err != nil || !token.Valid {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims, nil
	}
	
	return nil, errors.New("Failed to parse claims")
}
