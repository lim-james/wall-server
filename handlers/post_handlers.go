package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"wall-server/database"
	"wall-server/database/models"
)

type PostHandler struct {
	*Handler
}

func NewPostHandler(db *database.Database) *PostHandler {
	return &PostHandler{Handler: NewHandler(db)}
}

func (ph *PostHandler) ReadAllPostHandler(c *gin.Context) {
	posts, err := ph.DB.ReadAllPosts()

	if err != nil {
		ErrorResponse(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (ph *PostHandler) ReadAllPostsByUserIDHandler(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, errors.New("Invalid user_id"))
		return
	}

	posts, err := ph.DB.ReadAllPostsByUserID(userID)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (ph *PostHandler) CreatePostHandler(c *gin.Context) {
	var post models.Post

	if err := c.BindJSON(&post); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	post.UserID = c.MustGet("UserID").(int64)

	id, err := ph.DB.CreatePost(post)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	post.PostID = id
	c.JSON(http.StatusCreated, post)
}
