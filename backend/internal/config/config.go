package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	UploadDir  string
}

func MustLoad() Config {
	return Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "doc_hub"),
		DBUser:     getEnv("DB_USER", "doc_hub_user"),
		DBPassword: getEnv("DB_PASSWORD", "doc_hub_password"),
		UploadDir:  getEnv("UPLOAD_DIR", "uploads"),
	}
}

func (c Config) ServerAddress() string {
	return ":" + c.Port
}

func (c Config) DatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
