package service

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"go-teacher/backend/internal/domain"
)

var ErrInternal = errors.New("internal server error")

type documentRepository interface {
	Create(ctx context.Context, document domain.Document) (domain.Document, error)
	List(ctx context.Context) ([]domain.Document, error)
	GetByID(ctx context.Context, id int64) (domain.Document, error)
	Update(ctx context.Context, id int64, fileName string) (domain.Document, error)
	Delete(ctx context.Context, id int64) (domain.Document, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type fileStorage interface {
	Save(file multipart.File, fileName string) (string, int64, error)
	Delete(path string) error
}

type DocumentService struct {
	repository documentRepository
	storage    fileStorage
}

func NewDocumentService(repository documentRepository, storage fileStorage) *DocumentService {
	return &DocumentService{
		repository: repository,
		storage:    storage,
	}
}

func (s *DocumentService) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (domain.Document, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".pdf" && ext != ".docx" {
		return domain.Document{}, errors.New("only PDF and DOCX are allowed")
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = detectContentType(file)
	}

	storedName := buildStoredFileName(header.Filename)
	storagePath, size, err := s.storage.Save(file, storedName)
	if err != nil {
		return domain.Document{}, ErrInternal
	}

	document := domain.Document{
		OriginalName: header.Filename,
		StoredName:   storedName,
		ContentType:  contentType,
		Extension:    ext,
		Size:         size,
		StoragePath:  storagePath,
		Status:       "uploaded",
	}

	createdDocument, err := s.repository.Create(ctx, document)
	if err != nil {
		return domain.Document{}, ErrInternal
	}

	go s.processDocumentAsync(createdDocument.ID)

	return createdDocument, nil
}

func (s *DocumentService) List(ctx context.Context) ([]domain.Document, error) {
	documents, err := s.repository.List(ctx)
	if err != nil {
		return nil, ErrInternal
	}

	return documents, nil
}

func (s *DocumentService) GetByID(ctx context.Context, id int64) (domain.Document, error) {
	document, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.Document{}, mapRepositoryError(err)
	}

	return document, nil
}

func (s *DocumentService) Update(ctx context.Context, id int64, fileName string) (domain.Document, error) {
	if strings.TrimSpace(fileName) == "" {
		return domain.Document{}, errors.New("fileName is required")
	}

	document, err := s.repository.Update(ctx, id, fileName)
	if err != nil {
		return domain.Document{}, mapRepositoryError(err)
	}

	return document, nil
}

func (s *DocumentService) Delete(ctx context.Context, id int64) error {
	document, err := s.repository.Delete(ctx, id)
	if err != nil {
		return mapRepositoryError(err)
	}

	if err := s.storage.Delete(document.StoragePath); err != nil {
		return ErrInternal
	}

	return nil
}

func detectContentType(file multipart.File) string {
	buffer := make([]byte, 512)
	bytesRead, err := file.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return "application/octet-stream"
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "application/octet-stream"
	}

	return http.DetectContentType(buffer[:bytesRead])
}

func mapRepositoryError(err error) error {
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, repositoryErrNotFound()) {
		return ErrNotFound
	}

	return ErrInternal
}

func (s *DocumentService) processDocumentAsync(documentID int64) {
	time.Sleep(3 * time.Second)

	if err := s.repository.UpdateStatus(context.Background(), documentID, "processed"); err != nil {
		return
	}
}
