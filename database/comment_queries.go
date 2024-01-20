package database

import (
	"database/sql"
	"errors"
	"time"
	"wall-server/database/models"
)

const (
	selectAllCommentsByPostIDQuery = "SELECT comment_id, user_id, post_id, comment_text, creation_time, is_edited, IFNULL(last_edited_time, 'No Edit Time'), reply_id as last_edited_time FROM post_comments WHERE post_id = ?"
	selectCommentByIDQuery         = "SELECT comment_id, user_id, post_id, comment_text, creation_time, is_edited, IFNULL(last_edited_time, 'No Edit Time'), reply_id as last_edited_time FROM post_comments WHERE comment_id = ?"
	selectCommentAuthorByIDQuery   = "SELECT user_id FROM post_comments WHERE comment_id = ?"
	insertCommentQuery             = "INSERT INTO post_comments (post_id, user_id, comment_text, reply_id) VALUES (?, ?, ?, ?)"
	updateCommentQuery             = "UPDATE post_comments SET comment_text = ?, is_edited = TRUE, last_edited_time = CURRENT_TIMESTAMP WHERE comment_id = ?"
	selectLastEditedTimeQuery 		 = "SELECT last_edited_time FROM post_comments WHERE comment_id = ?"
	deleteCommentQuery             = "DELETE FROM post_comments WHERE comment_id = ?"
)

func (d *Database) ReadAllCommentsByPostID(postID int64) ([]models.Comment, error) {
	rows, err := d.DB.Query(selectAllCommentsByPostIDQuery, postID)
	if err != nil {
		return nil, HandleError(err)
	}
	defer rows.Close()

	comments := make([]models.Comment, 0)
	for rows.Next() {
		var comment models.Comment
		var creationTimeStr string
		var editedTimeStr string
		if err := rows.Scan(&comment.CommentID, &comment.UserID, &comment.PostID, &comment.Text, &creationTimeStr, &comment.IsEdited, &editedTimeStr, &comment.ReplyID); err != nil {
			return nil, HandleError(err)
		}

		comment.CreationTime, err = time.Parse("2006-01-02 15:04:05", creationTimeStr)
		if err != nil {
			return nil, HandleError(err)
		}

		if comment.IsEdited {
			comment.LastEditedTime, err = time.Parse("2006-01-02 15:04:05", editedTimeStr)
			if err != nil {
				return nil, HandleError(err)
			}
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (d *Database) ReadCommentByID(commentID int64, comment *models.Comment) error {
	var creationTimeStr string
	var editedTimeStr string

	err := d.DB.QueryRow(selectCommentByIDQuery, commentID).
		Scan(&comment.CommentID, &comment.UserID, &comment.PostID, &comment.Text, &creationTimeStr, &comment.IsEdited, &editedTimeStr, &comment.ReplyID)

	if err != nil {
		return HandleError(err)
	}

	comment.CreationTime, err = time.Parse("2006-01-02 15:04:05", creationTimeStr)
	if err != nil {
		return HandleError(err)
	}

	if comment.IsEdited {
		comment.LastEditedTime, err = time.Parse("2006-01-02 15:04:05", editedTimeStr)
		if err != nil {
			return HandleError(err)
		}
	}

	return nil
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
	if comment.ReplyID != 0 {
		var refComment models.Comment
		if err := d.ReadCommentByID(comment.ReplyID, &refComment); err != nil {
			return 0, HandleError(err)
		}
		if comment.PostID != refComment.PostID { 
			return 0, HandleError(errors.New("Comment reply_id does not match post_id"))
		}
	}

	var id int64

	err := d.withTransaction(func(tx *sql.Tx) error {
		result, err := tx.Exec(insertCommentQuery, comment.PostID, comment.UserID, comment.Text, comment.ReplyID)
		if err != nil {
			return err
		}

		id, err = result.LastInsertId()
		return err
	})

	return id, HandleError(err)
}

func (d *Database) EditComment(comment models.Comment) (time.Time, error) {
	var lastEditedTimeStr string
	err := d.withTransaction(func(tx *sql.Tx) error {
		_, err := tx.Exec(updateCommentQuery, comment.Text, comment.CommentID)

		if err != nil {
			return err
		}
		
		return tx.QueryRow(selectLastEditedTimeQuery, comment.CommentID).Scan(&lastEditedTimeStr)
	})

	if err != nil {
		return time.Time{}, HandleError(err)
	}
	
	lastEditTime, err := time.Parse("2006-01-02 15:04:05", lastEditedTimeStr)
	return lastEditTime, HandleError(err)
}

func (d *Database) DeleteComment(commentID int64) error {
	_, err := d.DB.Exec(deleteCommentQuery, commentID)
	return HandleError(err)
}
