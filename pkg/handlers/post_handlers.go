package handlers

import (
    "net/http"

    "wall-server/pkg/database"
    "wall-server/pkg/models"

    "github.com/gin-gonic/gin"
)

type Handler struct {
    DB *database.Database
}

func NewHandler(db *database.Database) *Handler {
    return &Handler{DB: db}
}

func (h *Handler) ReadAllPostHandler(c *gin.Context) {
    posts, err := h.DB.ReadAllPosts()

    if err != nil {
        c.IndentedJSON(http.StatusNotFound, err)
        return
    }

    c.IndentedJSON(http.StatusOK, posts)
}

func (h *Handler) CreatePostHandler(c *gin.Context) {
    var post models.Post

    if err := c.BindJSON(&post); err != nil {
        c.IndentedJSON(http.StatusBadRequest, err)
        return
    }

    id, err := h.DB.CreatePost(post)
    if err != nil {
        c.IndentedJSON(http.StatusBadRequest, err)
        return
    }

    post.ID = id
    c.IndentedJSON(http.StatusCreated, post)
}
