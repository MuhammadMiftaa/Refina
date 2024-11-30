package main

import (
	// "log"
	"os"

	"server/config"
	"server/interface/http/router"

	// "github.com/joho/godotenv"
)

func main() {
	// if _, err := os.Stat(".env"); err == nil {
	// 	if err := godotenv.Load(); err != nil {
	// 		log.Println("Error loading .env file")
	// 	}
	// }
	
	port := os.Getenv("PORT")

	db, err := config.SetupDatabase()
	if err != nil {
		panic(err)
	}

	redis := config.SetupRedisDatabase()

	r := router.SetupRouter(db, redis)
	r.Run(":" + port)
}
