package handlers

import (
	"net/http"

	"go-teacher/backend/internal/domain"
	httpresponse "go-teacher/backend/internal/http/response"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	response := domain.MessageResponse{
		Message: "Hello from Go API",
	}

	httpresponse.WriteJSON(w, http.StatusOK, response)
}
