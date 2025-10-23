package main

import (
	"armur-codescanner/internal/api"
	"armur-codescanner/internal/redis"
	"armur-codescanner/internal/worker"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"

	_ "armur-codescanner/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Armur Code Scanner API
// @version 1.0
// @description This is a code scanner service API.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /
func main() {
	// Default to binding on loopback for safer local development.
	bindAddr := os.Getenv("BIND_ADDR")
	if bindAddr == "" {
		bindAddr = "127.0.0.1"
	}

	router := gin.Default()
	go func() {
		if err := startAsynqWorker(); err != nil {
			log.Fatalf("Failed to start Asynq worker: %v", err)
		}
	}()
	api.RegisterRoutes(router)

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "4500"
	}
	addr := net.JoinHostPort(bindAddr, port)
	if err := router.Run(addr); err != nil {
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
