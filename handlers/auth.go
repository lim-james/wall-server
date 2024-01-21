package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"wall-server/database"
	"wall-server/database/models"
)

var jwtSecret = []byte("your-secret-key")

type AuthHandler struct {
	*Handler
}

func NewAuthHandler(db *database.Database) *AuthHandler {
	return &AuthHandler{Handler: NewHandler(db)}
}

func (ah *AuthHandler) LoginHandler(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	userCredentials, err := ah.DB.ReadUserByUsername(user.Username)

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if userCredentials == nil || userCredentials.PasswordHash != user.PasswordHash {
		ErrorResponse(c, http.StatusUnauthorized, errors.New("Invalid username or password"))
		return
	}

	token, err := ah.generateToken(userCredentials.UserID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ah *AuthHandler) SignupHandler(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	userID, err := ah.DB.CreateUser(user)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	token, err := ah.generateToken(userID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	user.UserID = userID
	c.JSON(http.StatusCreated, gin.H{"user": user, "token": token})
}

func (ah *AuthHandler) DeleteUserHandler(c *gin.Context) {
	username := c.Param("username")
	userID := c.MustGet("UserID").(int64)
	deletedUserID, err := ah.DB.ReadUserIDByUsername(username)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, errors.New("Invalid username"))
		return
	}

	if userID != deletedUserID {
		ErrorResponse(c, http.StatusUnauthorized, errors.New("You cannot delete other user"))
		return
	}

	
	if err := ah.DB.DeleteUser(userID); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (ah *AuthHandler) generateToken(userID int64) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
