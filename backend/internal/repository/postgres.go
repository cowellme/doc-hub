package repository

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenPostgres(dsn string) (*sql.DB, error) {
	return sql.Open("pgx", dsn)
}

func RunMigrations(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS documents (
			id BIGSERIAL PRIMARY KEY,
			original_name TEXT NOT NULL,
			stored_name TEXT NOT NULL,
			content_type TEXT NOT NULL,
			extension TEXT NOT NULL,
			size BIGINT NOT NULL,
			storage_path TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`

	_, err := db.Exec(query)
	return err
}
