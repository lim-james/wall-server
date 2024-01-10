package main

import (
    "database/sql"
    "log"
    "os"

    "wall-server/pkg/database"
    "wall-server/pkg/handlers"

    "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
    if err != nil {
        log.Fatal(err)
    }
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    database := database.NewDatabase(db)
    handler := handlers.NewHandler(database)

    r := gin.Default()

    r.GET("/", handler.ReadAllPostHandler)
    r.POST("/", handler.CreatePostHandler)

    r.Run("0.0.0.0:80")
}
