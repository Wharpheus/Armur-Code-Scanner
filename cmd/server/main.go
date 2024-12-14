package main

import (
	"armur-codescanner/internal/api"
	"armur-codescanner/internal/redis"
	"armur-codescanner/internal/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"log"
	"os"
)

func main() {
	router := gin.Default()
	go func() {
		if err := startAsynqWorker(); err != nil {
			log.Fatalf("Failed to start Asynq worker: %v", err)
		}
	}()
	api.RegisterRoutes(router)
	port := os.Getenv("APP_PORT")
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

func startAsynqWorker() error {
	server := asynq.NewServer(
		redis.RedisClientOptions(),
		asynq.Config{
			Concurrency: 10,
		},
	)

	mux := asynq.NewServeMux()
	mux.Handle("scan:repo", &worker.ScanTaskHandler{})

	// Start the Asynq server and process tasks
	return server.Start(mux)
}
