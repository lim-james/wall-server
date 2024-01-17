package models

import "time"

type Post struct {
	PostID       int64     `json:"post_id"`
	UserID       int64     `json:"user_id"`
	Title        string    `json:"title" binding:"required"`
	Body         string    `json:"body" binding:"required"`
	CreationTime time.Time `json:"creation_time"`
}
