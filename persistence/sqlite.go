package persistence

import (
	"fmt"
	"database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type DB struct {
	database* sql.DB
}

func (db* DB) Open(path string) (err error) {
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		fmt.Println("error opening the database", err)
		return err
	}
	db.database = database
	fmt.Println("database opened", database)
	return nil
}
