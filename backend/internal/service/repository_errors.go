package service

import "go-teacher/backend/internal/repository"

func repositoryErrNotFound() error {
	return repository.ErrDocumentNotFound
}
