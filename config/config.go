package config

import (
    "fmt"
    // "log"
    "os"

	"go-auth-service/pkg/logger"

    "github.com/joho/godotenv"
)

type Config struct {
    Port      string
    MongoURI  string
    JWTSecret string
    RedisURL  string
    NatsURL   string
}

func Load() (*Config, error) {
    
    if err := godotenv.Load(); err != nil {
		logger.Error(".env file not found")
	}

    cfg := &Config {
		Port:      os.Getenv("PORT"),
		MongoURI:  os.Getenv("MONGO_URI"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		RedisURL:  os.Getenv("REDIS_URL"),
		NatsURL:   os.Getenv("NATS_URL"),
	}

    if err:= validateConfig(cfg); err != nil {
        logger.Error("Config validation failed: " + err.Error())
        return nil, err
    }

    logger.Info("Configuration loaded successfully")
    return cfg, nil

}


func validateConfig(cfg *Config) error {
	required := map[string]string{
		"PORT":       cfg.Port,
		"MONGO_URI":  cfg.MongoURI,
		"JWT_SECRET": cfg.JWTSecret,
		"REDIS_URL":  cfg.RedisURL,
		"NATS_URL":   cfg.NatsURL,
	}

	for key, val := range required {
		if val == "" {
			return fmt.Errorf("Missing required env: %s", key)
		}
	}

	return nil
}