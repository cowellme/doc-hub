package repository

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenPostgres(dsn string) (*sql.DB, error) {
	return sql.Open("pgx", dsn)
}

func RunMigrations(db *sql.DB) error {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS documents (
			id BIGSERIAL PRIMARY KEY,
			original_name TEXT NOT NULL,
			stored_name TEXT NOT NULL,
			content_type TEXT NOT NULL,
			extension TEXT NOT NULL,
			size BIGINT NOT NULL,
			storage_path TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'uploaded',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`

	if _, err := db.Exec(createTableQuery); err != nil {
		return err
	}

	alterStatusQuery := `
		ALTER TABLE documents
		ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'uploaded'
	`

	_, err := db.Exec(alterStatusQuery)
	return err
}
