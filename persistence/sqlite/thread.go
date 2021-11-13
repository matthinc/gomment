package sqlite

import (
	"fmt"

	"github.com/matthinc/gomment/model"
)

func (db *DB) GetThreads() ([]model.Thread, error) {
	rows, err := db.database.Query("SELECT `thread_id`, `path` FROM `thread`")
	if err != nil {
		return nil, fmt.Errorf("failed to query database for thread: %w", err)
	}
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
