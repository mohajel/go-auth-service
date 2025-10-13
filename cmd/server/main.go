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
	// Load configuration
	cfg := config.Load()
	logger.Info("Config loaded")
	logger.Info("Starting Auth Service...")

	// --- MongoDB ---
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mgo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}
	db := mongoClient.Database("authdb")
	logger.Info("Connected to MongoDB")

	// --- Redis ---
	rdb := r.NewClient(&r.Options{Addr: cfg.RedisURL})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis connect error:", err)
	}
	logger.Info("Connected to Redis")

	// --- NATS ---
	nc := nats.Connect(cfg.NatsURL)
	if nc == nil {
		log.Fatal("Failed to connect to NATS")
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil || js == nil {
		log.Fatal("Failed to get JetStream context:", err)
	}

	// Create Stream
	nats.EnsureStream(js)

	// --- Subscriber ---
	subscriberReady := make(chan bool)
	go func() {
		nats.SubscribeUserRegistered(js)
		// Queue group example (optional)
		nats.SubscribeQueue(js, "worker-group-1")
		nats.SubscribeQueue(js, "worker-group-2")
		subscriberReady <- true
	}()
	<-subscriberReady

	// --- Publish test message ---
	nats.PublishUserRegistered(js, "user-123")

	// --- Repositories ---
	userRepo := user.NewMongoRepo(db)
	redisRepo := redis.NewRepository(rdb)

	// --- Token Manager ---
	tokenManager := token.NewManager(cfg.JWTSecret)

	// --- Auth Service & Handler ---
	authService := auth.NewService(userRepo, tokenManager, nc, redisRepo)
	authHandler := auth.NewHandler(authService)

	// --- Router ---
	router := gin.Default()
	authHandler.RegisterRoutes(router)

	logger.Info("Auth service running on port " + cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
