package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"wall-server/database/models"
)

func (ph *PostHandler) ReadAllCommentsHandler(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, errors.New("Invalid post_id"))
		return
	}

	comments, err := ph.DB.ReadAllCommentsByPostID(postID)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (ph *PostHandler) CreateCommentHandler(c *gin.Context) {
	userID := c.MustGet("UserID").(int64)

	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, errors.New("Invalid post_id"))
		return
	}

	var comment models.Comment

	if err := c.BindJSON(&comment); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	comment.UserID = userID
	comment.PostID = postID

	comment.CommentID, err = ph.DB.CreateComment(comment)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (ph *PostHandler) EditCommentHandler(c *gin.Context) {
	commentIDStr := c.Param("comment_id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, errors.New("Invalid comment_id"))
		return
	}

	userID := c.MustGet("UserID").(int64)

	var author int64
	if author, err = ph.DB.ReadCommentAuthorByID(commentID); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if author != userID {
		ErrorResponse(c, http.StatusBadRequest, errors.New("You are not the author of this comment"))
		return
	}

	var newComment models.Comment

	if err := c.BindJSON(&newComment); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	newComment.CommentID = commentID

	if err := ph.DB.EditComment(newComment); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, newComment)
}

func (ph *PostHandler) DeleteCommentHandler(c *gin.Context) {
	commentIDStr := c.Param("comment_id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, errors.New("Invalid comment_id"))
		return
	}

	userID := c.MustGet("UserID").(int64)

	var author int64
	if author, err = ph.DB.ReadCommentAuthorByID(commentID); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if author != userID {
		ErrorResponse(c, http.StatusBadRequest, errors.New("You are not the author of this comment"))
		return
	}

	if err := ph.DB.DeleteComment(commentID); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}
