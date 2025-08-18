package main

import (
	"context"
	"log"

	"server/config/db"
	"server/config/env"
	"server/config/queue"
	"server/internal/repository"
	"server/internal/service"
)

func init() {
	// Load environment variables and configuration
	if err := env.LoadByViper(); err != nil {
		log.Println("[ERROR] Failed to read JSON config file:", err)
		log.Println("[INFO] Loading environment variables from .env file")
		env.LoadNative()
	}
	db.SetupDatabase()      // Initialize the database connection and run migrations
	db.SetupRedisDatabase() // Initialize the Redis connection
	queue.SetupRabbitMQ()   // Initialize RabbitMQ connection
}

func main() {
	userRepo := repository.NewUsersRepository(db.DB)
	reportsRepo := repository.NewReportsRepository(db.DB)
	reportsService := service.NewReportsService(reportsRepo, userRepo)

	ctx := context.Background()
	if err := reportsService.UpdateUserReport(ctx); err != nil {
		log.Println("Failed to update user report: " + err.Error())
	}
}
