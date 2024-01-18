package database

import (
	"database/sql"
	"log"
	"time"
	"wall-server/database/models"
)

const (
	selectPostsQuery = "SELECT post_id, user_id, title, body, creation_time FROM posts"
	selectPostsByUserIDQuery  = "SELECT post_id, user_id, title, body, creation_time FROM posts WHERE user_id = ?"
	insertPostQuery  = "INSERT INTO posts (user_id, title, body) VALUES (?, ?, ?)"
)

func (d *Database) ReadAllPosts() ([]models.Post, error) {
	rows, err := d.DB.Query(selectPostsQuery)
	if err != nil {
		return nil, HandleError(err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var creationTimeStr string
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Body, &creationTimeStr); err != nil {
			log.Fatal(err)
			return nil, err
		}

		post.CreationTime, err = time.Parse("2006-01-02 15:04:05", creationTimeStr)
		if err != nil {
			return nil, HandleError(err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, HandleError(err)
	}

	return posts, nil
}

func (d *Database) ReadAllPostsByUserID(userID int64) ([]models.Post, error) {
	rows, err := d.DB.Query(selectPostsByUserIDQuery, userID)
	if err != nil {
		return nil, HandleError(err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var creationTimeStr string
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Body, &creationTimeStr); err != nil {
			log.Fatal(err)
			return nil, err
		}

		post.CreationTime, err = time.Parse("2006-01-02 15:04:05", creationTimeStr)
		if err != nil {
			return nil, HandleError(err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, HandleError(err)
	}

	return posts, nil
}
func (d *Database) CreatePost(post models.Post) (int64, error) {
	var id int64

	err := d.withTransaction(func(tx *sql.Tx) error {
		result, err := tx.Exec(insertPostQuery, post.UserID, post.Title, post.Body)
		if err != nil {
			return err
		}

		id, err = result.LastInsertId()
		return err
	})

	return id, HandleError(err)
}
