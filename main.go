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

type Post struct {
    ID int64 `json:"id"`
    Title string `json:"title"`
    Body string `json"body"`
}

func getAllHandler(c *gin.Context) {
    rows, err := db.Query("SELECT * FROM post")
    if err != nil {
        c.IndentedJSON(http.StatusNotFound, err)
        return;
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        if err := rows.Scan(&post.ID, &post.Title, &post.Body); err != nil {
            c.IndentedJSON(http.StatusNotFound, err)
            return
        }
        posts = append(posts, post)
    }

    if err := rows.Err(); err != nil {
        c.IndentedJSON(http.StatusNotFound, err)
        return
    }

    c.IndentedJSON(http.StatusOK, posts)
}

func newPost(post Post) (int64, error) {
    result, err := db.Exec("INSERT INTO post (title, body) VALUES (?, ?)", post.Title, post.Body)
    if err != nil {
        return -1, err
    }

    return result.LastInsertId()
}

func newPostHandler(c *gin.Context) {
    var post Post
    
    if err := c.BindJSON(&post); err != nil {
        c.IndentedJSON(http.StatusNotFound, err)
        return
    }

    id, err := newPost(post)
    if err != nil {
        c.IndentedJSON(http.StatusNotFound, err)
        return
    }

    post.ID = id;
    c.IndentedJSON(http.StatusCreated, post)
}

func main() {
    var err error
    db, err = sql.Open("mysql", os.Getenv("MYSQL_DSN"))
    if err != nil {
        log.Fatal(err)
    }
    
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    r := gin.Default()

    r.GET("/", getAllHandler)    
    r.POST("/", newPostHandler)    

    r.Run("0.0.0.0:80")
}
