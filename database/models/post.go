package models

import "time"

type Post struct {
	PostID         int64     `json:"post_id"`
	UserID         int64     `json:"user_id"`
	Title          string    `json:"title" binding:"required"`
	Body           string    `json:"body" binding:"required"`
	CreationTime   time.Time `json:"creation_time"`
	IsEdited       bool      `json:"is_edited"`
	LastEditedTime time.Time `json:"last_edited_time"`
}

type PostFormatted struct {
	PostID         int64     `json:"post_id"`
	Username       string    `json:"username"`
	Title          string    `json:"title" binding:"required"`
	Body           string    `json:"body" binding:"required"`
	CreationTime   time.Time `json:"creation_time"`
	IsEdited       bool      `json:"is_edited"`
	LastEditedTime time.Time `json:"last_edited_time"`
	LikeCount      int64     `json:"like_count"`
}

type PostDetailsFormatted struct {
	PostID         int64              `json:"post_id"`
	Username       string             `json:"username"`
	Title          string             `json:"title" binding:"required"`
	Body           string             `json:"body" binding:"required"`
	CreationTime   time.Time          `json:"creation_time"`
	IsEdited       bool               `json:"is_edited"`
	LastEditedTime time.Time          `json:"last_edited_time"`
	LikeCount      int64              `json:"like_count"`
	Comments       []CommentFormatted `json:"comments"`
}
