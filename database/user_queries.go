package database

import (
	"database/sql"
	"errors"
	"wall-server/database/models"
)

const (
	selectUserByUsernameQuery = "SELECT user_id, username, password_hash FROM users WHERE username = ?"
	insertUserQuery           = "INSERT INTO users (username, password_hash) VALUES (?, ?)"
)

func (d *Database) ReadUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := d.DB.QueryRow(selectUserByUsernameQuery, username).
		Scan(&user.UserID, &user.Username, &user.PasswordHash)

	if err != nil {
		return nil, HandleError(err)
	}

	return &user, nil
}

func (d *Database) CreateUser(user models.User) (int64, error) {
	var id int64
	err := d.withTransaction(func(tx *sql.Tx) error {
		if !IsUniqueUsername(tx, user.Username) {
			return errors.New("Username already exists")
		}

		result, err := tx.Exec(insertUserQuery, user.Username, user.PasswordHash)
		if err != nil {
			return err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		id = lastID
		return err
	})

	return id, HandleError(err)
}

func (d *Database) ReadUserIDByUsername(username string) (int64, error) {
	var userID int64
	err := d.DB.QueryRow("SELECT user_id FROM users WHERE username = ?", username).
		Scan(&userID)

	if err != nil {
		return 0, HandleError(err)
	}

	return userID, nil
}

func (d *Database) ReadUsernameByUserID(userID int64) (string, error) {
	var username string
	err := d.DB.QueryRow("SELECT username FROM users WHERE user_id = ?", userID).
		Scan(&username)

	if err != nil {
		return "", HandleError(err)
	}

	return username, nil
}

func (d *Database) DeleteUser(userID int64) error {
	return HandleError(d.withTransaction(func(tx *sql.Tx) error {
		_, err := tx.Exec("DELETE FROM users WHERE user_id = ?", userID)
		return err
	}))
}

func IsUniqueUsername(tx *sql.Tx, username string) bool {
	var count int
	err := tx.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		HandleError(err)
		return false
	}
	return count == 0
}
