package main

import (
	"server/config"
	"server/interface/router"
)

func main() {
	db, err := config.SetupDatabase()
	if err != nil {
		panic(err)
	}

	r := router.SetupRouter(db)
	r.Run(":8080")
}
