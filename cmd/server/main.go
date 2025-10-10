package main

import (
	"context"
	"log"
	"time"

	"go-auth-service/config"
	"go-auth-service/internal/auth"
	"go-auth-service/internal/nats"
	"go-auth-service/internal/redis"
	"go-auth-service/internal/token"
	"go-auth-service/internal/user"
	"go-auth-service/pkg/logger"

	"github.com/gin-gonic/gin"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	r "github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	logger.Info("Starting Auth Service...")


	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mgo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}
	db := mongoClient.Database("authdb")

	rdb := r.NewClient(&r.Options{
		Addr: cfg.RedisURL,
	})
	if err := rdb.Ping(ctx).Err(); err != nil { 
		log.Fatal("Redis connect error:", err)
	}

	nc := nats.Connect(cfg.NatsURL)

	userRepo := user.NewMongoRepo(db)
	redisRepo := redis.NewRepository(rdb)

	tokenManager := token.NewManager(cfg.JWTSecret)

	authService := auth.NewService(userRepo, tokenManager, nc, redisRepo)
	authHandler := auth.NewHandler(authService)

	router := gin.Default()
	authHandler.RegisterRoutes(router)

	logger.Info("Auth service running on port " + cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
