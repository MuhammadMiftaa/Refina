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
	host := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbName, dbPort)
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
