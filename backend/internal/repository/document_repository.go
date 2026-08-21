package repository

import (
	"context"
	"database/sql"
	"errors"

	"go-teacher/backend/internal/domain"
)

var ErrDocumentNotFound = errors.New("document not found")

type DocumentRepository struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(ctx context.Context, document domain.Document) (domain.Document, error) {
	query := `
		INSERT INTO documents (
			original_name,
			stored_name,
			content_type,
			extension,
			size,
			storage_path,
			status
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, status, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		document.OriginalName,
		document.StoredName,
		document.ContentType,
		document.Extension,
		document.Size,
		document.StoragePath,
		document.Status,
	).Scan(&document.ID, &document.Status, &document.CreatedAt)
	if err != nil {
		return domain.Document{}, err
	}

	return document, nil
}

func (r *DocumentRepository) List(ctx context.Context) ([]domain.Document, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, original_name, stored_name, content_type, extension, size, storage_path, status, created_at
		FROM documents
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []domain.Document
	for rows.Next() {
		var document domain.Document
		if err := rows.Scan(
			&document.ID,
			&document.OriginalName,
			&document.StoredName,
			&document.ContentType,
			&document.Extension,
			&document.Size,
			&document.StoragePath,
			&document.Status,
			&document.CreatedAt,
		); err != nil {
			return nil, err
		}

		documents = append(documents, document)
	}

	return documents, rows.Err()
}

func (r *DocumentRepository) GetByID(ctx context.Context, id int64) (domain.Document, error) {
	var document domain.Document

	err := r.db.QueryRowContext(ctx, `
		SELECT id, original_name, stored_name, content_type, extension, size, storage_path, status, created_at
		FROM documents
		WHERE id = $1
	`, id).Scan(
		&document.ID,
		&document.OriginalName,
		&document.StoredName,
		&document.ContentType,
		&document.Extension,
		&document.Size,
		&document.StoragePath,
		&document.Status,
		&document.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Document{}, ErrDocumentNotFound
	}
	if err != nil {
		return domain.Document{}, err
	}

	return document, nil
}

func (r *DocumentRepository) Update(ctx context.Context, id int64, fileName string) (domain.Document, error) {
	document, err := r.GetByID(ctx, id)
	if err != nil {
		return domain.Document{}, err
	}

	err = r.db.QueryRowContext(ctx, `
		UPDATE documents
		SET original_name = $2
		WHERE id = $1
		RETURNING id, original_name, stored_name, content_type, extension, size, storage_path, status, created_at
	`, id, fileName).Scan(
		&document.ID,
		&document.OriginalName,
		&document.StoredName,
		&document.ContentType,
		&document.Extension,
		&document.Size,
		&document.StoragePath,
		&document.Status,
		&document.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Document{}, ErrDocumentNotFound
	}
	if err != nil {
		return domain.Document{}, err
	}

	return document, nil
}

func (r *DocumentRepository) Delete(ctx context.Context, id int64) (domain.Document, error) {
	document, err := r.GetByID(ctx, id)
	if err != nil {
		return domain.Document{}, err
	}

	result, err := r.db.ExecContext(ctx, `DELETE FROM documents WHERE id = $1`, id)
	if err != nil {
		return domain.Document{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return domain.Document{}, err
	}
	if affected == 0 {
		return domain.Document{}, ErrDocumentNotFound
	}

	return document, nil
}

func (r *DocumentRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE documents
		SET status = $2
		WHERE id = $1
	`, id, status)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrDocumentNotFound
	}

	return nil
}
