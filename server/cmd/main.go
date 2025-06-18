package main

import (
	"log"
	"os"

	"server/db/config"
	"server/interface/http/router"

	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Error loading .env file")
		}
	}
	
	port := os.Getenv("PORT")
	log.Println("Starting server on port:", port)
	if port == "" {
		log.Println("No PORT environment variable set, using default port 8080")
		port = "8080" // Default port if not set
	}

	db, err := config.SetupDatabase()
	if err != nil {
		panic(err)
	}

	redis := config.SetupRedisDatabase()

	r := router.SetupRouter(db, redis)
	r.Run(":" + port)
}
