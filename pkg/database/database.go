package database

import (
    "database/sql"
    "wall-server/pkg/models"
)

type Database struct {
    DB *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
    return &Database{DB: db}
}

func (d *Database) ReadAllPosts() ([]models.Post, error) {
    rows, err := d.DB.Query("SELECT * FROM post")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []models.Post
    for rows.Next() {
        var post models.Post
        if err := rows.Scan(&post.ID, &post.Title, &post.Body); err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return posts, nil
}

func (d *Database) CreatePost(post models.Post) (int64, error) {
    result, err := d.DB.Exec("INSERT INTO post (title, body) VALUES (?, ?)", post.Title, post.Body)
    if err != nil {
        return 0, err
    }

    return result.LastInsertId()
}
