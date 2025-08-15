package main

import (
	dc "server/db/config"
	ec "server/env/config"
	qc "server/queue/config"
	"server/interface/http/router"
)

func init() {
	ec.LoadByViper() // Load environment variables and configuration
	dc.SetupDatabase() // Initialize the database connection and run migrations
	dc.SetupRedisDatabase() // Initialize the Redis connection
	qc.SetupRabbitMQ() // Initialize RabbitMQ connection
}

func main() {
	defer qc.Close() // Close RabbitMQ connection when the application exits

	r := router.SetupRouter() // Set up the HTTP router
	r.Run(":" + ec.Cfg.Server.Port)
}
