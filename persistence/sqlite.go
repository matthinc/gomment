package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/matthinc/gomment/model"
	_ "github.com/mattn/go-sqlite3"
)

type DBError struct {
	message string
}

func (err DBError) Error() string {
	return err.message
}

type DB struct {
	database *sql.DB
}

func New() DB {
	return DB{}
}

func (db *DB) Close() {
	db.database.Close()
}

const commentSelectFields = "`comment_id`, `parent_id`, `created_at`, `touched_at`, `num_children`, `author`, `text`"

func (db *DB) Open(path string) (err error) {
	// https://github.com/mattn/go-sqlite3/issues/377
	database, err := sql.Open("sqlite3", path+"?_foreign_keys=on")
	if err != nil {
		fmt.Println("error opening the database", err)
		return err
	}
	db.database = database
	return nil
}

// setup the database by creating all required tables
func (db *DB) Setup() (err error) {
	_, err1 := db.database.Exec("CREATE TABLE IF NOT EXISTS `thread` ( " +
		"`thread_id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, " +
		"`path` TEXT NOT NULL UNIQUE, " +
		"`num_total` INTEGER NOT NULL DEFAULT 1, " +
		"`num_root` INTEGER NOT NULL DEFAULT 1" +
		")")
	_, err2 := db.database.Exec("CREATE TABLE IF NOT EXISTS `comment` ( " +
		"`comment_id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, " +
		"`thread_id` INTEGER NOT NULL, " +
		"`parent_id` INTEGER DEFAULT NULL, " +
		"`num_children` INTEGER NOT NULL DEFAULT 0, " +
		"`verified` INTEGER NOT NULL DEFAULT 0, " +
		"`created_at` INTEGER NOT NULL, " +
		"`edited_at` INTEGER DEFAULT NULL, " +
		"`touched_at` INTEGER NOT NULL, " +
		"`author` TEXT NOT NULL, " +
		"`email` TEXT DEFAULT NULL, " +
		"`text` TEXT NOT NULL, " +
		"FOREIGN KEY(`thread_id`) REFERENCES `thread` (`thread_id`), " +
		"FOREIGN KEY(`parent_id`) REFERENCES `comment` (`comment_id`) " +
		")")

	if err1 != nil || err2 != nil {
		return DBError{"Unable to create DB"}
	}

	return nil
}

func (db *DB) CreateComment(commentCreation *model.CommentCreation, createdAt int64) (int64, error) {
	ctx := context.TODO()
	transaction, err := db.database.BeginTx(ctx, nil)

	numRoot := 0
	if commentCreation.ParentId == 0 {
		numRoot = 1
	}

	// create a new thread if it doesn't already exist
	res, err := transaction.ExecContext(ctx, "INSERT INTO `thread`(`path`) VALUES(?) ON CONFLICT(`path`) DO UPDATE SET "+
		"`num_total` = `num_total` + 1, "+
		"`num_root` = `num_root` + ?",
		commentCreation.ThreadPath,
		numRoot,
	)
	if err != nil {
		transaction.Rollback()
		return 0, fmt.Errorf("failed to insert thread: %w", err)
	}

	parentId := sql.NullInt64{
		Int64: int64(commentCreation.ParentId),
		Valid: commentCreation.ParentId != 0,
	}

	res, err = transaction.Exec(
		"INSERT INTO `comment` (`thread_id`, `parent_id`, `created_at`, `touched_at`, `author`, `email`, `text`) "+
			"SELECT t1.`thread_id`, ?, ?, ?, ?, ?, ? FROM `thread` t1 WHERE t1.`path` = ?",
		parentId,
		createdAt,
		createdAt,
		commentCreation.Author,
		commentCreation.Email,
		commentCreation.Text,
		commentCreation.ThreadPath,
	)
	if err != nil {
		transaction.Rollback()
		return 0, fmt.Errorf("failed to insert comment: %w", err)
	}

	commentId, err := res.LastInsertId()
	if err != nil {
		transaction.Rollback()
		return 0, fmt.Errorf("failed to retrieve id of recently inserted comment: %w", err)
	}

	// set the touched_at date of all (recursive) parent comments to the created_at date of the leaf comment
	res, err = transaction.Exec(
		"UPDATE `comment` SET `touched_at` = c1.`created_at` FROM (SELECT `created_at` FROM `comment` WHERE `comment_id` = ?) AS c1 WHERE `comment_id` IN ("+
			"WITH RECURSIVE `parents`(`parent_id`) AS (SELECT ? UNION SELECT `comment`.`parent_id` FROM `comment`,`parents` where `comment_id` = `parents`.`parent_id`)"+
			"SELECT `parent_id` from parents"+
			")",
		commentId,
		commentId,
	)
	if err != nil {
		transaction.Rollback()
		return 0, fmt.Errorf("failed to update the `touched_at` times for parent comments: %w", err)
	}

	// increment the num_children variable for the parent
	if parentId.Valid {
		res, err = transaction.Exec(
			"UPDATE `comment` SET `num_children` = `num_children` + 1 WHERE `comment_id` = ?;",
			parentId,
		)
		if err != nil {
			transaction.Rollback()
			return 0, fmt.Errorf("failed to increment the `num_children` value for parent comment: %w", err)
		}
	}

	err = transaction.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit the comment insertion transaction: %w", err)
	}

	return commentId, nil
}

func (db *DB) parseCommentsQuery(rows *sql.Rows) ([]model.Comment, error) {
	response := make([]model.Comment, 0)

	var (
		id          int
		parent      sql.NullInt64
		createdAt   int64
		touchedAt   int64
		numChildren int
		author      string
		text        string
	)

	for rows.Next() {
		err := rows.Scan(&id, &parent, &createdAt, &touchedAt, &numChildren, &author, &text)
		if err != nil {
			return nil, fmt.Errorf("failed to scan result row: %w", err)
		}

		parentId := 0
		if parent.Valid {
			parentId = int(parent.Int64)
		}

		comment := model.Comment{
			Id:          id,
			Author:      author,
			Email:       "",
			Text:        text,
			ParentId:    parentId,
			CreatedAt:   createdAt,
			TouchedAt:   touchedAt,
			NumChildren: numChildren,
		}
		response = append(response, comment)
	}

	return response, nil
}

func (db *DB) QueryCommentsById(thread int) ([]model.Comment, error) {
	rows, err := db.database.Query("SELECT "+commentSelectFields+" FROM `comment` where thread_id = ? ORDER BY created_at DESC", thread)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	return db.parseCommentsQuery(rows)
}

func (db *DB) GetNewestCommentsByPath(path string, limit int) ([]model.Comment, ThreadMetaInfo, error) {
	rows, err := db.database.Query(
		"SELECT `thread_id`, `num_total`, `num_root` FROM `thread` WHERE `path` = ?",
		path,
	)
	defer rows.Close()
	if err != nil {
		return nil, ThreadMetaInfo{}, fmt.Errorf("failed to query database for thread: %w", err)
	}

	var (
		threadId int
		numTotal int
		numRoot  int
	)

	if !rows.Next() {
		// the thread does not exist
		return []model.Comment{}, ThreadMetaInfo{}, nil
	}

	err = rows.Scan(&threadId, &numTotal, &numRoot)
	if err != nil {
		return nil, ThreadMetaInfo{}, fmt.Errorf("failed to scan fields for thread query: %w", err)
	}

	rows, err = db.database.Query(
		"SELECT "+commentSelectFields+" FROM `comment` WHERE `thread_id` = ? ORDER BY `touched_at` DESC, `created_at` ASC LIMIT ?",
		threadId,
		limit,
	)
	defer rows.Close()

	if err != nil {
		return nil, ThreadMetaInfo{}, fmt.Errorf("failed to query database for comments: %w", err)
	}

	comments, err := db.parseCommentsQuery(rows)

	return comments, ThreadMetaInfo{
		NumTotal: numTotal,
		NumRoot:  numRoot,
	}, err
}

func (db *DB) GetThreads() ([]model.Thread, error) {
	rows, _ := db.database.Query("SELECT `thread_id`, `path` FROM `thread`")
	defer rows.Close()

	response := make([]model.Thread, 0)

	var id int
	var path string
	for rows.Next() {
		if err := rows.Scan(&id, &path); err != nil {
			return nil, fmt.Errorf("failed to scan result row: %w", err)
		}

		response = append(response, model.Thread{
			Id:   id,
			Path: path,
		})
	}

	return response, nil
}
