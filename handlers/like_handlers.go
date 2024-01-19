package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ph *PostHandler) LikePostHandler(c *gin.Context) {
	userID := c.MustGet("UserID").(int64)

	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, errors.New("Invalid post_id"))
		return
	}

	if has, _ := ph.DB.HasLikedPost(userID, postID); has {
		ErrorResponse(c, http.StatusBadRequest, errors.New("You have already liked this post"))
		return
	}

	if err := ph.DB.LikePost(userID, postID); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, errors.New("You have already liked this post"))
		return
	}

	// Fetch and return the total number of likes for the post
	totalLikes, err := ph.DB.GetTotalLikesForPost(postID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, errors.New("Failed to fetch total likes"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post liked successfully", "total_likes": totalLikes})
}

func (ph *PostHandler) UnlikePostHandler(c *gin.Context) {
	userID := c.MustGet("UserID").(int64)
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, errors.New("Invalid post_id"))
		return
	}

	if has, _ := ph.DB.HasLikedPost(userID, postID); !has {
		ErrorResponse(c, http.StatusBadRequest, errors.New("You have not liked this post"))
		return
	}

	if err := ph.DB.UnlikePost(userID, postID); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, errors.New("Failed to unlike the post"))
		return
	}

	// Fetch and return the total number of likes for the post
	totalLikes, err := ph.DB.GetTotalLikesForPost(postID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, errors.New("Failed to fetch total likes"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post liked successfully", "total_likes": totalLikes})
}
