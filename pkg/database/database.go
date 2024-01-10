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
	tx, err := d.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) 
		} else if err != nil {
			_ = tx.Rollback() 
		} else {
			err = tx.Commit() 
		}
	}()

	result, err := tx.Exec("INSERT INTO post (title, body) VALUES (?, ?)", post.Title, post.Body)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}