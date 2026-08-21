package domain

import "time"

type UploadResponse struct {
	ID             int64     `json:"id"`
	FileName       string    `json:"fileName"`
	StoredFileName string    `json:"storedFileName"`
	ContentType    string    `json:"contentType"`
	Size           int64     `json:"size"`
	StoragePath    string    `json:"storagePath"`
	Message        string    `json:"message"`
	UploadedAt     time.Time `json:"uploadedAt"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Document struct {
	ID           int64
	OriginalName string
	StoredName   string
	ContentType  string
	Extension    string
	Size         int64
	StoragePath  string
	CreatedAt    time.Time
}

type UpdateDocumentRequest struct {
	FileName string `json:"fileName"`
}
