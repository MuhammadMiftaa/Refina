package main

import (
	"server/config/db"
	"server/config/env"
	"server/config/log"
	"server/config/queue"
	"server/interface/http/router"
)

func init() {
	log.SetupLogger() // Initialize the logger configuration
	log.Info("Starting Refina API...")

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
	db.SetupRedisDatabase() // Initialize the Redis connection
	log.Info("Setup Redis Connection Success")

	log.Info("Setup RabbitMQ Connection Start")
	queue.SetupRabbitMQ() // Initialize RabbitMQ connection
	log.Info("Setup RabbitMQ Connection Success")
}

func main() {
	defer log.Info("Refina API stopped")
	defer queue.Close() // Close RabbitMQ connection when the application exits

	r := router.SetupRouter() // Set up the HTTP router
	r.Run(":" + env.Cfg.Server.Port)
	log.Info("Starting HTTP server on port " + env.Cfg.Server.Port)
}
