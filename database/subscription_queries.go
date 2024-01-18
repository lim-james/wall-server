package database

import (
	"database/sql"
	"wall-server/database/models"
)

const (
	selectSubscriptionQuery 			= "SELECT post_id FROM subscriptions WHERE subscriber_id = ?"
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
		var postID int64
		if err := rows.Scan(&postID); err != nil {
			return nil, HandleError(err)
		}

		var post models.Post
		if err := d.ReadPostByID(postID, &post); err != nil {
			return nil, HandleError(err)
		}

		posts = append(posts, post)
	}	

	if err := rows.Err(); err != nil {
		return nil, HandleError(err)
	}

	return posts, nil
}