package sqlite

import (
	"fmt"

	"go.uber.org/zap"
)

type migrationFunc func(db *DB) error

var migrations = map[uint]migrationFunc{
	0: create_db, // right when db has been created
	1: migration_1_to_2,
}

func (db *DB) getSchemaVersion() (uint, error) {
	rows, err := db.database.Query("PRAGMA user_version;")
	if err != nil {
		return 0, fmt.Errorf("error while retrieving DB schema version: %w", err)
	}
	defer rows.Close()

	var ret uint
	rows.Next()
	err = rows.Scan(&ret)
	if err != nil {
		return 0, fmt.Errorf("error while parsing DB schema version: %w", err)
	}

	return ret, nil
}

func (db *DB) setSchemaVersion(version uint) error {
	_, err := db.database.Exec(fmt.Sprintf("PRAGMA user_version = %d;", version))
	if err != nil {
		return fmt.Errorf("error while setting DB schema version: %w", err)
	}

	return nil
}

func (db *DB) doMigrations() error {
	schemaVersion, err := db.getSchemaVersion()
	if err != nil {
		return err
	}

	for migrationFn, ok := migrations[schemaVersion]; ok; migrationFn, ok = migrations[schemaVersion] {
		err = migrationFn(db)
		if err != nil {
			return fmt.Errorf("migration from DB schema version %d failed: %w", schemaVersion, err)
		}

		newSchemaVersion, err := db.getSchemaVersion()
		if err != nil {
			return err
		}

		if newSchemaVersion <= schemaVersion {
			return fmt.Errorf("migration from DB schema version %d did not increase DB schema version", schemaVersion)
		}

		zap.L().Sugar().Infof("DB migration changed schema version %d to %d", schemaVersion, newSchemaVersion)
		schemaVersion = newSchemaVersion
	}

	return nil
}

func create_db(db *DB) error {
	zap.L().Info("running DB migration 'create_db'")

	_, err := db.database.Exec("CREATE TABLE `thread` ( " +
		"`thread_id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, " +
		"`path` TEXT NOT NULL UNIQUE, " +
		"`num_total` INTEGER NOT NULL DEFAULT 1, " +
		"`num_root` INTEGER NOT NULL DEFAULT 1" +
		")")
	if err != nil {
		return fmt.Errorf("creation of table 'thread' failed: %w", err)
	}

	_, err = db.database.Exec("CREATE TABLE `comment` ( " +
		"`comment_id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, " +
		"`thread_id` INTEGER NOT NULL, " +
		"`parent_id` INTEGER DEFAULT NULL, " +
		"`num_children` INTEGER NOT NULL DEFAULT 0, " +
		"`depth_level` INTEGER NOT NULL DEFAULT 0, " +
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
	if err != nil {
		return fmt.Errorf("creation of table 'comment' failed: %w", err)
	}

	_, err = db.database.Exec("CREATE INDEX `thread_index` ON `comment` (`thread_id` ASC);")
	if err != nil {
		return fmt.Errorf("creation of index 'thread_index' failed: %w", err)
	}

	return db.setSchemaVersion(2)
}

func migration_1_to_2(db *DB) error {
	zap.L().Info("running DB migration 'migration_1_to_2'")

	_, err := db.database.Exec("CREATE INDEX `thread_index` ON `comment` (`thread_id` ASC);")
	if err != nil {
		return fmt.Errorf("creation of index 'thread_index' failed: %w", err)
	}

	return db.setSchemaVersion(2)
}
