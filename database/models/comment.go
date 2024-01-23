package models

import (
	"time"
)

type Comment struct {
	CommentID      int64     `json:"comment_id"`
	PostID         int64     `json:"post_id"`
	UserID         int64     `json:"user_id"`
	Text           string    `json:"text" binding:"required"`
	CreationTime   time.Time `json:"creation_time"`
	ReplyID        int64     `json:"reply_id"`
	IsEdited       bool      `json:"is_edited"`
	LastEditedTime time.Time `json:"last_edited_time"`
}

type CommentFormatted struct {
	CommentID      int64     `json:"comment_id"`
	PostID         int64     `json:"post_id"`
	Username       string    `json:"username"`
	Text           string    `json:"text"` 
	CreationTime   time.Time `json:"creation_time"`
	ReplyID        int64     `json:"reply_id"`
	IsEdited       bool      `json:"is_edited"`
	LastEditedTime time.Time `json:"last_edited_time"`
}