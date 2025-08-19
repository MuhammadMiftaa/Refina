package db

import (
	"fmt"

	"server/config/env"
	"server/config/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", env.Cfg.Database.DBHost, env.Cfg.Database.DBUser, env.Cfg.Database.DBPassword, env.Cfg.Database.DBName, env.Cfg.Database.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	DB = db
}
