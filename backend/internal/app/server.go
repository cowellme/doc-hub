package app

import (
	"database/sql"
	"log"
	"net/http"

	"go-teacher/backend/internal/config"
	"go-teacher/backend/internal/http/handlers"
	"go-teacher/backend/internal/http/middleware"
	"go-teacher/backend/internal/repository"
	"go-teacher/backend/internal/service"
	"go-teacher/backend/internal/storage"
)

func NewServer(cfg config.Config) (*http.Server, func()) {
	db := openDatabase(cfg)
	documentRepository := repository.NewDocumentRepository(db)
	documentService := service.NewDocumentService(documentRepository, storage.NewLocalFileStorage(cfg.UploadDir))

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", handlers.Hello)
	mux.HandleFunc("/documents", handlers.Documents(documentService))
	mux.HandleFunc("/documents/", handlers.DocumentByID(documentService))
	mux.HandleFunc("/upload", handlers.Documents(documentService))
	mux.HandleFunc("/openapi.yaml", handlers.OpenAPI)
	mux.HandleFunc("/swagger", handlers.SwaggerUI)

	server := &http.Server{
		Addr:    cfg.ServerAddress(),
		Handler: middleware.CORS(mux),
	}

	return server, func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}
}

func openDatabase(cfg config.Config) *sql.DB {
	db, err := repository.OpenPostgres(cfg.DatabaseDSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := repository.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	return db
}
