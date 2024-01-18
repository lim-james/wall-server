package models

import (
	"time"
)

type Comment struct {
	CommentID    int64     `json:"comment_id"`
	PostID       int64     `json:"post_id"`
	UserID       int64     `json:"user_id"`
	Text         string    `json:"text" binding:"required"`
	CreationTime time.Time `json:"creation_time"`
}
