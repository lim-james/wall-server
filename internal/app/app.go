package app

import (
	"database/sql"
	"log"
	"os"

	"wall-server/pkg/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func setupDatabase() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func Run() error {
	db := setupDatabase()
	defer db.Close()

	r := gin.Default()

	handlers.RegisterHandlers(r, db)

	return r.Run("0.0.0.0:80")
}
