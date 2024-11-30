package config

import (
	"fmt"
	"os"

	"server/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
)

func SetupDatabase() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=refina port=5432 sslmode=disable TimeZone=Asia/Jakarta", user, password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	db.AutoMigrate(&entity.Users{}, &entity.Transactions{})

	return db, err
}

func SetupRedisDatabase() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return rdb
}
