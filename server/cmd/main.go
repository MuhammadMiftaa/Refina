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

	r := router.SetupRouter(db)
	r.Run(":8080")
}
