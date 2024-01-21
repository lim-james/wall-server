package app

import (
	"database/sql"
	"os"

	"wall-server/database"
	"wall-server/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Server struct {
	DB     *sql.DB
	Router *gin.Engine
}

func NewServer() (*Server, error) {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	r := gin.Default()

	server := &Server{
		DB:     db,
		Router: r,
	}
	return server, nil
}

func (s *Server) Close() {
	s.DB.Close()
}

func (s *Server) Run() error {
	defer s.Close()

	database := database.NewDatabase(s.DB)

	authHandler := handlers.NewAuthHandler(database)
	postHandler := handlers.NewPostHandler(database)

	api := s.Router.Group("/api")
	{
		api.GET("/", postHandler.ReadAllPostHandler)
		api.GET("/p/:post_id/comments", postHandler.ReadAllCommentsHandler)

		post := api.Group("/p")
		{
			post.Use(handlers.AuthMiddleware())
			post.POST("/", postHandler.CreatePostHandler)
			post.PUT("/:post_id/", postHandler.EditPostHandler)
			post.DELETE("/:post_id/", postHandler.DeletePostHandler)
			post.POST("/:post_id/like", postHandler.LikePostHandler)
			post.POST("/:post_id/unlike", postHandler.UnlikePostHandler)
			post.POST("/:post_id/subscribe", postHandler.SubscribePostHandler)
			post.POST("/:post_id/unsubscribe", postHandler.UnsubscribePostHandler)
			post.POST("/:post_id/comments", postHandler.CreateCommentHandler)
			post.PUT("/:post_id/comments/:comment_id", postHandler.EditCommentHandler)
			post.DELETE("/:post_id/comments/:comment_id", postHandler.DeleteCommentHandler)
		}

		// Auth routes
		auth := api.Group("/u")
		{
			auth.POST("/signup", authHandler.SignupHandler)
			auth.POST("/login", authHandler.LoginHandler)
			auth.GET("/:username", postHandler.ReadAllPostsByUserIDHandler)
			auth.GET("/:username/subscriptions", handlers.AuthMiddleware(), postHandler.ReadAllSubscribedPostsHandler)
			auth.DELETE("/:username", handlers.AuthMiddleware(), authHandler.DeleteUserHandler)
		}
	}

	return s.Router.Run("0.0.0.0:80")
}
