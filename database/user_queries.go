package database

import (
	"database/sql"
	"errors"
	"wall-server/database/models"
)

const (
	selectUserByUsernameQuery = "SELECT user_id, username, password_hash FROM users WHERE username = ?"
	insertUserQuery = "INSERT INTO users (username, password_hash) VALUES (?, ?)"
)

func (d *Database) ReadUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := d.DB.QueryRow(selectUserByUsernameQuery, username).
		Scan(&user.UserID, &user.Username, &user.PasswordHash)

	if err != nil {
		// fmt.Println(err)
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

func IsUniqueUsername(tx *sql.Tx, username string) bool {
	var count int
	err := tx.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		HandleError(err)
		return false 
	}
	return count == 0
}
