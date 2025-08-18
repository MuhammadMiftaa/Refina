package main

import (
	"log"

	"server/config/db"
	"server/config/env"
	"server/config/queue"
	"server/interface/http/router"
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
	// defer queue.Close() // Close RabbitMQ connection when the application exits

	r := router.SetupRouter() // Set up the HTTP router
	r.Run(":" + env.Cfg.Server.Port)
}
