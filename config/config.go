package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	JWTSecret     string
	JWTAccessTTL  int
	JWTRefreshTTL int

	MongoURI string
	MongoDB  string
}

var App *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	App = &Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		AppPort:       getEnv("APP_PORT", "8080"),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTAccessTTL:  getEnvInt("JWT_ACCESS_TTL", 3600),
		JWTRefreshTTL: getEnvInt("JWT_REFRESH_TTL", 604800),
		MongoURI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDB:       getEnv("MONGODB_DB", "todolist"),
	}

	if App.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
