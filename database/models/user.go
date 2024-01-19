package models

type User struct {
	UserID       int64  `json:"user_id"`
	Username     string `json:"username" 			binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
}
