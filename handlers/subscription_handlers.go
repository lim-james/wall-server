package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ph *PostHandler) SubscribePostHandler(c *gin.Context) {
	userID := c.MustGet("UserID").(int64)
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post_id"})
		return
	}

	if has, _ := ph.DB.HasSubscribedPost(userID, postID); has {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already subscribed to this post"})
		return
	}

	if err := ph.DB.SubscribePost(userID, postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe the post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post subscribed successfully"})
}

func (ph *PostHandler) UnsubscribePostHandler(c *gin.Context) {
	userID := c.MustGet("UserID").(int64)
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post_id"})
		return
	}

	if has, _ := ph.DB.HasSubscribedPost(userID, postID); !has {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You are already unsubscribed to this post"})
		return
	}

	if err := ph.DB.UnsubscribePost(userID, postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe the post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post unsubscribed successfully"})
}

func (ph *PostHandler) ReadAllSubscribedPostsHandler(c *gin.Context) {
	userID := c.MustGet("UserID").(int64)

	posts, err := ph.DB.ReadAllSubscribedPosts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscribed posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}
