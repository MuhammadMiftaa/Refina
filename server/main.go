package main

import (
	"server/config"
	"server/interface/http/router"
)

func main() {
	db, err := config.SetupDatabase()
	if err != nil {
		panic(err)
	}

	redis := config.SetupRedisDatabase()

	r := router.SetupRouter(db, redis)
	r.Run(":8080")
}
