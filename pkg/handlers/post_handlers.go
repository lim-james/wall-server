package handlers

import (
	"database/sql"
	"net/http"

	"wall-server/pkg/models"

	"github.com/gin-gonic/gin"
)

// Handler struct to encapsulate dependencies
type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

func (h *Handler) ReadAllPostQuery() ([]models.Post, error) {
	rows, err := h.DB.Query("SELECT * FROM post")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Body); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (h *Handler) ReadAllPostHandler(c *gin.Context) {
	posts, err := h.ReadAllPostQuery()

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, posts)
}

func (h *Handler) CreatePostQuery(post models.Post) (int64, error) {
	result, err := h.DB.Exec("INSERT INTO post (title, body) VALUES (?, ?)", post.Title, post.Body)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (h *Handler) CreatePostHandler(c *gin.Context) {
	var post models.Post

	if err := c.BindJSON(&post); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	id, err := h.CreatePostQuery(post)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	post.ID = id
	c.IndentedJSON(http.StatusCreated, post)
}

func RegisterHandlers(r *gin.Engine, db *sql.DB) {
	handler := NewHandler(db)

	r.GET("/", handler.ReadAllPostHandler)
	r.POST("/", handler.CreatePostHandler)
}
