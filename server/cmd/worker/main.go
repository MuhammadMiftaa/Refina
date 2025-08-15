package main

import (
	"context"

	"server/db/config"
	dc "server/db/config"
	ec "server/env/config"
	"server/internal/repository"
	"server/internal/service"
	qc "server/queue/config"
)

func init() {
	ec.LoadByViper()               // Load environment variables and configuration
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
		panic("Failed to update user report: " + err.Error())
	}
}
