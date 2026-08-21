package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-teacher/backend/internal/domain"
	httpresponse "go-teacher/backend/internal/http/response"
	"go-teacher/backend/internal/service"
)

func Documents(documentService *service.DocumentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			documents, err := documentService.List(r.Context())
			if err != nil {
				httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}

			httpresponse.WriteJSON(w, http.StatusOK, toUploadResponses(documents, "documents loaded"))
		case http.MethodPost:
			file, header, err := r.FormFile("file")
			if err != nil {
				httpresponse.WriteError(w, http.StatusBadRequest, "file is required")
				return
			}
			defer file.Close()

			document, err := documentService.Upload(r.Context(), file, header)
			if err != nil {
				writeServiceError(w, err)
				return
			}

			httpresponse.WriteJSON(w, http.StatusCreated, toUploadResponse(document, "file uploaded successfully"))
		default:
			httpresponse.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func DocumentByID(documentService *service.DocumentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseDocumentID(r.URL.Path)
		if err != nil {
			httpresponse.WriteError(w, http.StatusBadRequest, "invalid document id")
			return
		}

		switch r.Method {
		case http.MethodGet:
			document, err := documentService.GetByID(r.Context(), id)
			if err != nil {
				writeServiceError(w, err)
				return
			}

			httpresponse.WriteJSON(w, http.StatusOK, toUploadResponse(document, "document loaded"))
		case http.MethodPut:
			var request domain.UpdateDocumentRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				httpresponse.WriteError(w, http.StatusBadRequest, "invalid request body")
				return
			}

			document, err := documentService.Update(r.Context(), id, request.FileName)
			if err != nil {
				writeServiceError(w, err)
				return
			}

			httpresponse.WriteJSON(w, http.StatusOK, toUploadResponse(document, "document updated"))
		case http.MethodDelete:
			if err := documentService.Delete(r.Context(), id); err != nil {
				writeServiceError(w, err)
				return
			}

			httpresponse.WriteJSON(w, http.StatusOK, domain.MessageResponse{Message: "document deleted"})
		default:
			httpresponse.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func parseDocumentID(path string) (int64, error) {
	idPart := strings.TrimPrefix(path, "/documents/")
	return strconv.ParseInt(idPart, 10, 64)
}

func toUploadResponses(documents []domain.Document, message string) []domain.UploadResponse {
	result := make([]domain.UploadResponse, 0, len(documents))
	for _, document := range documents {
		result = append(result, toUploadResponse(document, message))
	}

	return result
}

func toUploadResponse(document domain.Document, message string) domain.UploadResponse {
	return domain.UploadResponse{
		ID:             document.ID,
		FileName:       document.OriginalName,
		StoredFileName: document.StoredName,
		ContentType:    document.ContentType,
		Size:           document.Size,
		StoragePath:    document.StoragePath,
		Status:         document.Status,
		Message:        message,
		UploadedAt:     document.CreatedAt,
	}
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrNotFound:
		httpresponse.WriteError(w, http.StatusNotFound, err.Error())
	case service.ErrInternal:
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
	default:
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
	}
}
