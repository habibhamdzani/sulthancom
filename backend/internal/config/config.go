package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	MongoURI     string
	MongoDB      string
	RedisAddr    string
	RedisPassword string
	JWTSecret    string
	JWTExpiry    int
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:          getEnv("PORT", "8080"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:       getEnv("MONGO_DB", "sulthancom_pos"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		JWTSecret:     getEnv("JWT_SECRET", "sulthancom-secret-key-change-in-production"),
		JWTExpiry:     24,
	}

	if cfg.JWTSecret == "sulthancom-secret-key-change-in-production" {
		log.Println("WARNING: Using default JWT secret. Change JWT_SECRET in production!")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
