package main

import (
	"log"

	ec "server/env/config"
	"server/interface/http/router"
)

func init() {
	// Load environment variables and configuration
	if err := ec.LoadByViper(); err != nil {
		log.Println("[ERROR] Failed to read JSON config file:", err)
		log.Println("[INFO] Loading environment variables from .env file")
		ec.LoadNative()
	}
	// dc.SetupDatabase()      // Initialize the database connection and run migrations
	// dc.SetupRedisDatabase() // Initialize the Redis connection
	// qc.SetupRabbitMQ()      // Initialize RabbitMQ connection
}

func main() {
	// defer qc.Close() // Close RabbitMQ connection when the application exits

	r := router.SetupRouter() // Set up the HTTP router
	r.Run(":" + ec.Cfg.Server.Port)
}
