package app

import (
	"database/sql"
	"os"

	"wall-server/database"
	"wall-server/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func setupDatabase() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		database.HandleError(err)
	}
	if err := db.Ping(); err != nil {
		database.HandleError(err)
	}
	return db
}

func Run() error {
	db := setupDatabase()
	database := database.NewDatabase(db)
	defer db.Close()

	r := gin.Default()

	authHandler := handlers.NewAuthHandler(database)
	postHandler := handlers.NewPostHandler(database)

	api := r.Group("/api") 
	{
		api.GET("/", postHandler.ReadAllPostHandler)

		post := api.Group("/p")
		{
			post.Use(handlers.AuthMiddleware())
			post.POST("/", postHandler.CreatePostHandler)
			post.POST("/:post_id/like", postHandler.LikePostHandler)
			post.POST("/:post_id/unlike", postHandler.UnlikePostHandler)
		}

		// Auth routes
		auth := api.Group("/u")
		{
			auth.POST("/signup", authHandler.SignupHandler) 
			auth.POST("/login", authHandler.LoginHandler)
		}
	}

	return r.Run("0.0.0.0:80")
}
