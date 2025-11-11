package main

import (
	"context"
	"log"
	"time"

	"go-auth-service/config"
	"go-auth-service/internal/auth"
	appnats "go-auth-service/internal/nats"
	"go-auth-service/internal/redis"
	"go-auth-service/internal/token"
	"go-auth-service/internal/user"
	"go-auth-service/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	r "github.com/redis/go-redis/v9"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// --- Load Config ---
	cfg := config.Load()
	logger.Info("✅ Config loaded")
	logger.Info("Starting Auth Service...")

	// --- MongoDB ---
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mgo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}
	db := mongoClient.Database("authdb")
	logger.Info("✅ Connected to MongoDB")

	// --- Redis ---
	rdb := r.NewClient(&r.Options{Addr: cfg.RedisURL})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis connect error:", err)
	}
	logger.Info("✅ Connected to Redis")

	// --- NATS ---
	nc := appnats.Connect(cfg.NatsURL)
	if nc == nil {
		log.Fatal("❌ Failed to connect to NATS")
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil || js == nil {
		log.Fatal("❌ Failed to get JetStream context:", err)
	}

	if err := appnats.EnsureStream(js); err != nil {
		log.Fatal("❌ Failed to ensure NATS stream:", err)
	}

	// --- Subscribers ---
	go func() {
		// Wait a bit for stream to be fully ready
		time.Sleep(1 * time.Second)
		appnats.SubscribeUserRegistered(js)
		appnats.SubscribeQueue(js, "worker-group-1")
		appnats.SubscribeQueue(js, "worker-group-2")
	}()

	// --- Test Publish ---
	// Wait a bit for subscribers to be ready
	time.Sleep(2 * time.Second)
	if err := appnats.PublishUserRegistered(js, "user-123", "test@example.com", "Test User"); err != nil {
		logger.Error("Failed to publish test message: " + err.Error())
	}

	// --- Repositories ---
	userRepo := user.NewMongoRepo(db)
	redisRepo := redis.NewRepository(rdb)

	// --- Token Manager ---
	tokenManager := token.NewManager(cfg.JWTSecret)

	// --- Auth Service & Handler ---
	authService := auth.NewService(userRepo, tokenManager, nc, redisRepo)
	authHandler := auth.NewHandler(authService)

	// --- Gin Router ---
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	authHandler.RegisterRoutes(router)

	logger.Info("Auth service running on port " + cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
