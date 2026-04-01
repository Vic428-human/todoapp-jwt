package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	Port          string
	JWTSecret     string
	GCSBucketName string
}

func Load() (*Config, error) {
	// 本機有 .env 就讀，沒有也不要中止
	wd, _ := os.Getwd()
	log.Println("current working directory:", wd)
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env: %v", err)
	} else {
		log.Println(".env loaded successfully")
	}

	cfg := &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		Port:          os.Getenv("PORT"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		GCSBucketName: os.Getenv("GCS_BUCKET_NAME"),
	}

	// 可選：本機預設值
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	log.Printf("DatabaseURL: %q", cfg.DatabaseURL)
	if cfg.DatabaseURL == "" {
		log.Println("warning: DATABASE_URL is empty")
	}

	return cfg, nil
}
