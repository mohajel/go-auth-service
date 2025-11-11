package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	MongoURI  string
	JWTSecret string
	RedisURL  string
	NatsURL   string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		Port:               getEnv("PORT", "8080"),
		MongoURI:           getEnv("MONGO_URI", "mongodb://localhost:27017"),
		JWTSecret:          getEnv("JWT_SECRET", "secret"),
		RedisURL:           getEnv("REDIS_URL", "localhost:6379"),
		NatsURL:            getEnv("NATS_URL", "nats://localhost:4222"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
	}

	log.Println("Config loaded")
	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
