package main

import (
	"context"

	"server/config/db"
	"server/config/env"
	"server/config/log"
	"server/config/queue"
	"server/config/redis"
	"server/internal/repository"
	"server/internal/service"
)

func init() {
	log.SetupLogger() // Initialize the logger configuration

	var err error
	var missing []string
	if missing, err = env.LoadByViper(); err != nil {
		log.Error("Failed to read JSON config file:" + err.Error())
		log.Info("Switch loading environment variables to .env file")
		if missing, err = env.LoadNative(); err != nil {
			log.Log.Fatalf("Failed to load environment variables: %v", err)
		}
		log.Info("Environment variables by .env file loaded successfully")
	} else {
		log.Info("Environment variables by Viper loaded successfully")
	}

	if len(missing) > 0 {
		for _, envVar := range missing {
			log.Warn("Missing environment variable: " + envVar)
		}
	}

	log.Info("Setup Database Connection Start")
	db.SetupDatabase() // Initialize the database connection and run migrations
	log.Info("Setup Database Connection Success")

	log.Info("Setup Redis Connection Start")
	redis.SetupRedisDatabase() // Initialize the Redis connection
	log.Info("Setup Redis Connection Success")

	log.Info("Setup RabbitMQ Connection Start")
	queue.SetupRabbitMQ() // Initialize RabbitMQ connection
	log.Info("Setup RabbitMQ Connection Success")

	log.Info("Starting Refina worker...")
}

func main() {
	defer log.Info("Refina worker stopped")
	defer queue.Close() // Close RabbitMQ connection when the application exits

	userRepo := repository.NewUsersRepository(db.DB)
	reportsRepo := repository.NewReportsRepository(db.DB)
	reportsService := service.NewReportsService(reportsRepo, userRepo)

	ctx := context.Background()
	if err := reportsService.UpdateUserReport(ctx); err != nil {
		log.Error("Failed to update user report: " + err.Error())
	}

	log.Info("Refina worker started successfully")
}
