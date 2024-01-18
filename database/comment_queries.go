package database

import (
	"database/sql"
	"time"
	"wall-server/database/models"
)

const (
	selectAllCommentsByPostIDQuery  = "SELECT comment_id, user_id, post_id, comment_text, creation_time FROM post_comments WHERE post_id = ?"
	selectCommentAuthorByIDQuery  = "SELECT user_id FROM post_comments WHERE comment_id = ?"
	insertCommentQuery  = "INSERT INTO post_comments (post_id, user_id, comment_text) VALUES (?, ?, ?)"
	updateCommentQuery  = "UPDATE post_comments SET comment_text = ? WHERE comment_id = ?"
	deleteCommentQuery  = "DELETE FROM post_comments WHERE comment_id = ?"
)

func (d *Database) ReadAllCommentsByPostID(postID int64) ([]models.Comment, error) {
	rows, err := d.DB.Query(selectAllCommentsByPostIDQuery, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]models.Comment, 0)
	for rows.Next() {
		var comment models.Comment
		var creationTimeStr string
		if err := rows.Scan(&comment.CommentID, &comment.UserID, &comment.PostID, &comment.Text, &creationTimeStr); err != nil {
			return nil, HandleError(err)
		}

		comment.CreationTime, err = time.Parse("2006-01-02 15:04:05", creationTimeStr)
		if err != nil {
			return nil, HandleError(err)
		}
		
		comments = append(comments, comment)
	}

	return comments, nil
}

func (d *Database) ReadCommentAuthorByID(commentID int64) (int64, error) {
	var userID int64

	err := d.DB.QueryRow(selectCommentAuthorByIDQuery, commentID).
		Scan(&userID)

	if err != nil {
		return 0, HandleError(err)
	}

	return userID, nil
}

func (d *Database) CreateComment(comment models.Comment) (int64, error) {
	var id int64

	err := d.withTransaction(func(tx *sql.Tx) error {
		result, err := tx.Exec(insertCommentQuery, comment.PostID, comment.UserID, comment.Text)
		if err != nil {
			return err
		}

		id, err = result.LastInsertId()
		return err
	})

	return id, HandleError(err)
}

func (d *Database) EditComment(comment models.Comment) error {
	return HandleError(d.withTransaction(func(tx *sql.Tx) error {
		_, err := tx.Exec(updateCommentQuery, comment.Text, comment.CommentID)
		return err
	}))
}

func (d *Database) DeleteComment(commentID int64) error {
	_, err := d.DB.Exec(deleteCommentQuery, commentID)
	return HandleError(err)
}