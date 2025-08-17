package config

import (
	"fmt"
	"log"

	"server/env/config"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)

func SetupDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", config.Cfg.Database.DBHost, config.Cfg.Database.DBUser, config.Cfg.Database.DBPassword, config.Cfg.Database.DBName, config.Cfg.Database.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
		panic(err)
	}

	DB = db
}

func SetupRedisDatabase() {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.Cfg.Redis.RHost, config.Cfg.Redis.RPort),
	})

	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		log.Fatalf("Gagal terhubung ke Redis: %v", err)
		panic(err)
	}

	RDB = rdb
}
