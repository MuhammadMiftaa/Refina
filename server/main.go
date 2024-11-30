package main

import (
	"os"

	"server/config"
	"server/interface/http/router"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	port := os.Getenv("PORT")

	db, err := config.SetupDatabase()
	if err != nil {
		panic(err)
	}

	redis := config.SetupRedisDatabase()

	r := router.SetupRouter(db, redis)
	r.Run(":" + port)
}
