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

func HandleError(err error) error {
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (d *Database) withTransaction(fn func(*sql.Tx) error) error {
	tx, err := d.DB.Begin()
	if err != nil {
		return err
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

	return fn(tx)
}
