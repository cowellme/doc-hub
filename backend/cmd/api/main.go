package main

import (
	"log"

	"go-teacher/backend/internal/app"
	"go-teacher/backend/internal/config"
)

func main() {
	cfg := config.MustLoad()
	server, cleanup := app.NewServer(cfg)
	defer cleanup()

	log.Printf("server started on %s", cfg.ServerAddress())
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
