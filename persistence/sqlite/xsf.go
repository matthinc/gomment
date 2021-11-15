package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/matthinc/gomment/model"
	"github.com/matthinc/gomment/persistence"
)

func getVariableDepthXsfQuery(depth int, order string, limit int) string {
	if depth <= 0 {
		return "SELECT " + commentSelectFields + " FROM `comment` WHERE `thread_id` = ? AND `parent_id` IS NULL ORDER BY `created_at` " + order + " LIMIT " + fmt.Sprint(limit)
	}
	if depth == 1 {
		limit2 := int(float32(limit) * float32(0.25))
		limit1 := limit - limit2

		return "WITH anchor AS(" +
			"SELECT c1.`comment_id` FROM `comment` c1 WHERE c1.`thread_id` = ? AND c1.`parent_id` IS NULL ORDER BY c1.`created_at` " + order + " LIMIT " + fmt.Sprint(limit1) +
			"), lvl1 AS (" +
			"SELECT c1.`comment_id` FROM `comment` c1, anchor WHERE c1.`parent_id` = anchor.`comment_id` ORDER BY c1.`created_at` " + order + " LIMIT " + fmt.Sprint(limit2) +
			")" +
			"SELECT " + commentSelectFields + " FROM `comment` WHERE `comment_id` IN (SELECT * FROM anchor UNION SELECT * FROM lvl1) ORDER BY `created_at` " + order
	}

	limit3 := int(float32(limit) * float32(0.10))
	limit2 := int(float32(limit) * float32(0.20))
	limit1 := limit - limit3 - limit2

	return "WITH anchor AS(" +
		"SELECT c1.`comment_id` FROM `comment` c1 WHERE c1.`thread_id` = ? AND c1.`parent_id` IS NULL ORDER BY c1.`created_at` " + order + " LIMIT " + fmt.Sprint(limit1) +
		"), lvl1 AS (" +
		"SELECT c1.`comment_id` FROM `comment` c1, anchor WHERE c1.`parent_id` = anchor.`comment_id` ORDER BY c1.`created_at` " + order + " LIMIT " + fmt.Sprint(limit2) +
		"), lvl2 AS (" +
		"SELECT c1.`comment_id` FROM `comment` c1, lvl1 WHERE c1.`parent_id` = lvl1.`comment_id` ORDER BY c1.`created_at` " + order + " LIMIT " + fmt.Sprint(limit3) +
		")" +
		"SELECT " + commentSelectFields + " FROM `comment` WHERE `comment_id` IN (SELECT * FROM anchor UNION SELECT * FROM lvl1 UNION SELECT * FROM lvl2) ORDER BY `created_at` " + order
}

func (db *DB) getCommentsXsf(path string, maxDepth int, limit int, asc bool) ([]model.Comment, persistence.ThreadMetaInfo, error) {
	metaInfo, err := db.GetThreadMetaInfo(path)
	if err != nil {
		return nil, persistence.ThreadMetaInfo{}, fmt.Errorf("failed to retreive meta information about thread: %w", err)
	}

	order := "DESC"
	if asc {
		order = "ASC"
	}

	rows, err := db.database.Query(
		getVariableDepthXsfQuery(maxDepth, order, limit),
		metaInfo.ThreadId,
	)
	if err != nil {
		return nil, persistence.ThreadMetaInfo{}, fmt.Errorf("failed to query database for comments: %w", err)
	}
	defer rows.Close()

	comments, err := db.parseCommentsQuery(rows)

	return comments, metaInfo, err
}

func (db *DB) GetCommentsNsf(path string, maxDepth int, limit int) ([]model.Comment, persistence.ThreadMetaInfo, error) {
	return db.getCommentsXsf(path, maxDepth, limit, false)
}

func (db *DB) GetCommentsOsf(path string, maxDepth int, limit int) ([]model.Comment, persistence.ThreadMetaInfo, error) {
	return db.getCommentsXsf(path, maxDepth, limit, true)
}

func (db *DB) getMoreCommentsXsf(threadId int64, parentId int64, newestCreatedAt int64, limit int, asc bool) ([]model.Comment, error) {
	var err error
	var rows *sql.Rows

	order := "DESC"
	if asc {
		order = "ASC"
	}

	if parentId == 0 {
		rows, err = db.database.Query(
			"SELECT "+commentSelectFields+" FROM `comment` "+
				"WHERE `thread_id` = ? "+
				"AND `parent_id` IS NULL "+
				"AND `created_at` <= ? "+
				" ORDER BY `created_at` "+order+" LIMIT ?",
			threadId,
			newestCreatedAt,
			limit,
		)
	} else {
		rows, err = db.database.Query(
			"SELECT "+commentSelectFields+" FROM `comment` "+
				"WHERE `thread_id` = ? "+
				"AND `parent_id` = ? "+
				"AND `created_at` <= ? "+
				" ORDER BY `created_at` "+order+" LIMIT ?",
			threadId,
			parentId,
			newestCreatedAt,
			limit,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query database for thread: %w", err)
	}
	defer rows.Close()

	comments, err := db.parseCommentsQuery(rows)
	if err != nil {
		return nil, fmt.Errorf("failed scan comment result rows: %w", err)
	}

	return comments, nil
}

func (db *DB) GetMoreCommentsNsf(threadId int64, parentId int64, newestCreatedAt int64, limit int) ([]model.Comment, error) {
	return db.getMoreCommentsXsf(threadId, parentId, newestCreatedAt, limit, false)
}

func (db *DB) GetMoreCommentsOsf(threadId int64, parentId int64, newestCreatedAt int64, limit int) ([]model.Comment, error) {
	return db.getMoreCommentsXsf(threadId, parentId, newestCreatedAt, limit, true)
}
