package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/matthinc/gomment/model"
	"github.com/matthinc/gomment/persistence"
)

func (db *DB) GetCommentsNbf(path string, maxDepth int, limit int) ([]model.Comment, persistence.ThreadMetaInfo, error) {
	metaInfo, err := db.GetThreadMetaInfo(path)
	if err != nil {
		return nil, persistence.ThreadMetaInfo{}, fmt.Errorf("failed to retreive meta information about thread: %w", err)
	}

	rows, err := db.database.Query(
		"SELECT "+commentSelectFields+" FROM `comment` WHERE `thread_id` = ? AND `depth_level` < ? ORDER BY `touched_at` DESC, `created_at` ASC LIMIT ?",
		metaInfo.ThreadId,
		maxDepth,
		limit,
	)
	if err != nil {
		return nil, persistence.ThreadMetaInfo{}, fmt.Errorf("failed to query database for comments: %w", err)
	}
	defer rows.Close()

	comments, err := db.parseCommentsQuery(rows)

	return comments, metaInfo, err
}

func (db *DB) GetMoreCommentsNbf(threadId int64, parentId int64, newestCreatedAt int64, excludeIds []int64, limit int) ([]model.Comment, error) {
	// precondition: excludeIds must be ordered asc
	var lastId int64 = 0
	for _, id := range excludeIds {
		if id <= lastId {
			return nil, fmt.Errorf("excludeIds are not in ascending order")
		}
		lastId = id
	}

	var err error
	var rows *sql.Rows

	if parentId == 0 {
		rows, err = db.database.Query(
			"SELECT "+commentSelectFields+" FROM `comment` "+
				"WHERE `thread_id` = ? "+
				"AND `parent_id` IS NULL "+
				"AND `created_at` <= ? "+
				" ORDER BY `created_at` DESC LIMIT ?",
			threadId,
			newestCreatedAt,
			limit+len(excludeIds),
		)
	} else {
		rows, err = db.database.Query(
			"SELECT "+commentSelectFields+" FROM `comment` "+
				"WHERE `thread_id` = ? "+
				"AND `parent_id` = ? "+
				"AND `created_at` <= ? "+
				" ORDER BY `created_at` DESC LIMIT ?",
			threadId,
			parentId,
			newestCreatedAt,
			limit+len(excludeIds),
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

	ret := make([]model.Comment, 0)

	// 'ORDER BY `created_at` DESC' implies 'ORDER BY `comment_id` DESC'
	for _, comment := range comments {
		// remove all exclude id's until the id is <= the current comment
		for len(excludeIds) > 0 && excludeIds[len(excludeIds)-1] > int64(comment.Id) {
			excludeIds = excludeIds[:len(excludeIds)-1]
		}

		// keep the comment if it is not to be excluded
		if len(excludeIds) == 0 || int64(comment.Id) != excludeIds[len(excludeIds)-1] {
			ret = append(ret, comment)
		}
	}

	return ret, nil
}
