package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"wall-server/database"
)

// CustomClaims represents the custom claims you might have in your JWT token
type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

type Handler struct {
	DB *database.Database
}

func NewHandler(db *database.Database) *Handler {
	return &Handler{DB: db}
}

func ErrorResponse(c *gin.Context, statusCode int, err error) {
	c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
}
