package database

import (
	"database/sql"
)

const (
	selectLikesCountQuery    = "SELECT COUNT(*) AS total_likes FROM post_likes WHERE post_id = ?"
	selectLikeExistsQuery    = "SELECT 1 FROM post_likes WHERE user_id = ? AND post_id = ? LIMIT 1"
	insertLikeQuery          = "INSERT INTO post_likes (user_id, post_id) VALUES (?, ?)"
	deleteLikeQuery          = "DELETE FROM post_likes WHERE user_id = ? AND post_id = ?"
)

func (d *Database) GetTotalLikesForPost(postID int64) (int, error) {
	var total int
	err := d.DB.QueryRow(selectLikesCountQuery, postID).Scan(&total)
	if err != nil {
		return 0, HandleError(err)
	}

	return total, nil
}

func (d *Database) HasLikedPost(userID, postID int64) (bool, error) {
	var exists int
	err := d.DB.QueryRow(selectLikeExistsQuery, userID, postID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, HandleError(err)
	}
	return true, nil
}

func (d *Database) LikePost(userID, postID int64) error {
	_, err := d.DB.Exec(insertLikeQuery, userID, postID)
	return HandleError(err)
}

func (d *Database) UnlikePost(userID, postID int64) error {
	_, err := d.DB.Exec(deleteLikeQuery, userID, postID)
	return HandleError(err)
}