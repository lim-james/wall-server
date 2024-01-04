package main

import (
    "os"
    "log"

    "net/http"
    "database/sql"

    "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("mysql", os.Getenv("MYSQL_DSN"))
    if err != nil {
        log.Fatal(err)
    }
    
    ping := "We are LIVE!"
    if err := db.Ping(); err != nil {
        log.Fatal(err)
        ping = "Failed to connect"
    }

    r := gin.Default()

    r.GET("/", func (c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": ping})
    })    

    r.Run("0.0.0.0:80")
}
