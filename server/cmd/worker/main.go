package main

import (
	"context"
	"log"

	"server/db/config"
	dc "server/db/config"
	ec "server/env/config"
	"server/internal/repository"
	"server/internal/service"
	qc "server/queue/config"
)

func init() {
	// Load environment variables and configuration
	if err := ec.LoadByViper(); err != nil {
		log.Println("[ERROR] Failed to read JSON config file:", err)
		log.Println("[INFO] Loading environment variables from .env file")
		ec.LoadNative()
	}
	dc.SetupDatabase()      // Initialize the database connection and run migrations
	dc.SetupRedisDatabase() // Initialize the Redis connection
	qc.SetupRabbitMQ()      // Initialize RabbitMQ connection
}

func main() {
	userRepo := repository.NewUsersRepository(config.DB)
	reportsRepo := repository.NewReportsRepository(config.DB)
	reportsService := service.NewReportsService(reportsRepo, userRepo)

	ctx := context.Background()
	if err := reportsService.UpdateUserReport(ctx); err != nil {
		log.Println("Failed to update user report: " + err.Error())
	}
}
