package database

import (
	"database/sql"
	"time"
	"wall-server/database/models"
)

const (
	selectSubscriptionQuery       = "SELECT p.post_id, p.user_id, p.title, p.body, p.creation_time, p.is_edited, IFNULL(p.last_edited_time, 'No Edit Time') as last_edited_time FROM posts p INNER JOIN subscriptions s ON p.post_id = s.post_id WHERE s.subscriber_id = ?"
	selectSubscriptionExistsQuery = "SELECT 1 FROM subscriptions WHERE subscriber_id = ? AND post_id = ? LIMIT 1"
	insertSubscriptionQuery       = "INSERT INTO subscriptions (subscriber_id, post_id) VALUES (?, ?)"
	deleteSubscriptionQuery       = "DELETE FROM subscriptions WHERE subscriber_id = ? AND post_id = ?"
)

func (d *Database) HasSubscribedPost(userID, postID int64) (bool, error) {
	var exists int
	err := d.DB.QueryRow(selectSubscriptionExistsQuery, userID, postID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, HandleError(err)
	}
	return true, nil
}

func (d *Database) SubscribePost(userID, postID int64) error {
	_, err := d.DB.Exec(insertSubscriptionQuery, userID, postID)
	return HandleError(err)
}

func (d *Database) UnsubscribePost(userID, postID int64) error {
	_, err := d.DB.Exec(deleteSubscriptionQuery, userID, postID)
	return HandleError(err)
}

func (d *Database) ReadAllSubscribedPosts(userID int64) ([]models.Post, error) {
	rows, err := d.DB.Query(selectSubscriptionQuery, userID)
	if err != nil {
		return nil, HandleError(err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var creationTimeStr string
		var editedTimeStr string
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Body, &creationTimeStr, &post.IsEdited, &editedTimeStr); err != nil {
			return nil, HandleError(err)
		}

		post.CreationTime, err = time.Parse("2006-01-02 15:04:05", creationTimeStr)
		if err != nil {
			return nil, HandleError(err)
		}

		if post.IsEdited {
			post.LastEditedTime, err = time.Parse("2006-01-02 15:04:05", editedTimeStr)
			if err != nil {
				return nil, HandleError(err)
			}
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, HandleError(err)
	}

	return posts, nil
}
