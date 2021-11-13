package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/matthinc/gomment/model"
)

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
		"INSERT INTO `comment` (`thread_id`, `parent_id`, `depth_level`, `created_at`, `touched_at`, `author`, `email`, `text`) "+
			"SELECT t1.`thread_id`, ?, (SELECT ifnull(MAX(p.`depth_level`), -1) + 1 FROM `comment` p WHERE p.`comment_id` = ?), ?, ?, ?, ?, ? FROM `thread` t1 WHERE t1.`path` = ?",
		parentId,
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
