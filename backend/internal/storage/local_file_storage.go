package storage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalFileStorage struct {
	baseDir string
}

func NewLocalFileStorage(baseDir string) *LocalFileStorage {
	return &LocalFileStorage{baseDir: baseDir}
}

func (s *LocalFileStorage) Save(file multipart.File, fileName string) (string, int64, error) {
	if err := os.MkdirAll(s.baseDir, 0o755); err != nil {
		return "", 0, err
	}

	fullPath := filepath.Join(s.baseDir, fileName)
	destination, err := os.Create(fullPath)
	if err != nil {
		return "", 0, err
	}
	defer destination.Close()

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", 0, err
	}

	size, err := io.Copy(destination, file)
	if err != nil {
		return "", 0, err
	}

	return fullPath, size, nil
}

func (s *LocalFileStorage) Delete(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
