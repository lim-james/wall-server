package database

import (
	"database/sql"
	"log"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{DB: db}
}
