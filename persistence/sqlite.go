package persistence

import (
	"fmt"
	"database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type DBError struct {
	message string
}

func (err DBError) Error() string {
	return err.message
}

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

func (db* DB) Setup() (err error) {
	_, err1 := db.database.Exec("CREATE TABLE `thread` ( `thread_id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE )")
	_, err2 :=db.database.Exec("CREATE TABLE `comment` ( `comment_id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, `thread_id` INTEGER NOT NULL, `parent_id` INTEGER DEFAULT NULL, `verified` INTEGER NOT NULL DEFAULT 0, `created_at` INTEGER NOT NULL, `edited_at` INTEGER DEFAULT NULL, `author` TEXT, `email` TEXT, `text` TEXT )")

	if err1 != nil || err2 != nil {
		return DBError {"Unable to create DB"}
	}

	return nil
}
