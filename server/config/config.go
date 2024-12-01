package config

import (
	"fmt"
	"log"
	"os"

	"server/internal/entity"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var mode = func() string {
	m := os.Getenv("MODE")
	if m == "" {
		m = "development"
	}
	return m
}()

func SetupDatabase() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	host := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	var dsn string
	if mode == "development" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbName, dbPort)
	} else if mode == "production" {
		dsn = os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Fatal("DATABASE_URL tidak ditemukan di environment variables")
		}
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	if err := db.AutoMigrate(&entity.Users{}, &entity.Transactions{}); err != nil {
		log.Fatalf("Error saat melakukan migrasi: %v", err)
	}

	return db, nil
}

func SetupRedisDatabase() *redis.Client {
	var rdb *redis.Client
	if mode == "development" {
		rdb = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
	} else if mode == "production" {
		redisAddr := os.Getenv("REDIS_URL")
		if redisAddr == "" {
			log.Fatal("REDIS_URL tidak ditemukan di environment variables")
		}

		opt, err := redis.ParseURL(redisAddr)
		if err != nil {
			log.Fatalf("Error parsing REDIS_URL: %v", err)
		}

		rdb = redis.NewClient(opt)
	}

	return rdb
}
