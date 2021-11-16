package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/matthinc/gomment/model"
	"github.com/matthinc/gomment/persistence"

	_ "github.com/mattn/go-sqlite3"
)

type CommentRow struct {
	CommentId   int64
	ThreadId    int64
	ParentId    sql.NullInt64
	NumChildren int
	DepthLevel  int
	Verified    bool
	CreatedAt   int64
	EditedAt    sql.NullInt64
	TouchedAt   int64
	Author      string
	Email       sql.NullString
	Text        string
}

const commentSelectFields = "`comment_id`, `parent_id`, `created_at`, `touched_at`, `num_children`, `author`, `text`"

type DB struct {
	database *sql.DB
}

func New() DB {
	return DB{}
}

func (db *DB) Close() {
	db.database.Close()
}

func (db *DB) Open(path string) (err error) {
	// https://github.com/mattn/go-sqlite3/issues/377
	database, err := sql.Open("sqlite3", path+"?_foreign_keys=on")
	if err != nil {
		return fmt.Errorf("error opening the database: %w", err)
	}
	db.database = database
	return nil
}

func (db *DB) Setup() error {
	err := db.doMigrations()
	if err != nil {
		return fmt.Errorf("initial DB setup failed in migration: %w", err)
	}

	return nil
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

func (db *DB) GetCommentRow(commentId int64) (CommentRow, error) {
	rows, err := db.database.Query("SELECT * FROM `comment` WHERE `comment_id` = ?", commentId)
	if err != nil {
		return CommentRow{}, err
	}

	var ret CommentRow

	rows.Next()
	err = rows.Scan(
		&ret.CommentId,
		&ret.ThreadId,
		&ret.ParentId,
		&ret.NumChildren,
		&ret.DepthLevel,
		&ret.Verified,
		&ret.CreatedAt,
		&ret.EditedAt,
		&ret.TouchedAt,
		&ret.Author,
		&ret.Email,
		&ret.Text,
	)
	if err != nil {
		return CommentRow{}, err
	}

	return ret, nil
}

func (db *DB) GetThreadMetaInfo(path string) (persistence.ThreadMetaInfo, error) {
	rows, err := db.database.Query(
		"SELECT `thread_id`, `num_total`, `num_root` FROM `thread` WHERE `path` = ?",
		path,
	)
	if err != nil {
		return persistence.ThreadMetaInfo{}, fmt.Errorf("failed to query database for thread: %w", err)
	}
	defer rows.Close()

	var ret persistence.ThreadMetaInfo

	if !rows.Next() {
		// the thread does not exist
		return persistence.ThreadMetaInfo{}, nil
	}

	err = rows.Scan(&ret.ThreadId, &ret.NumTotal, &ret.NumRoot)
	if err != nil {
		return persistence.ThreadMetaInfo{}, fmt.Errorf("failed to scan fields for thread query: %w", err)
	}

	return ret, nil
}
