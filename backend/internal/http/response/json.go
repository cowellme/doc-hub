package response

import (
	"encoding/json"
	"net/http"

	"go-teacher/backend/internal/domain"
)

func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	WriteJSON(w, statusCode, domain.ErrorResponse{Message: message})
}
