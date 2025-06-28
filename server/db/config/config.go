package config

import (
	"fmt"
	"log"
	"os"

	"server/db/prepopulate"
	"server/db/views"
	"server/internal/types/entity"

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

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	if err := db.AutoMigrate(&entity.Users{}, &entity.WalletTypes{}, &entity.Wallets{}, &entity.Categories{}, &entity.Transactions{}, &entity.Attachments{}, &entity.InvestmentTypes{}, &entity.Investments{}); err != nil {
		log.Fatalf("Error saat melakukan migrasi: %v", err)
	}

	PrePopulate(db)
	CreateView(db)

	return db, nil
}

func SetupRedisDatabase() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})

	return rdb
}

func PrePopulate(db *gorm.DB) {
	prepopulate.WalletTypesSeeder(db)
	prepopulate.CategoryTypeSeeder(db)
	prepopulate.InvestmentTypesSeeder(db)
}

func CreateView(db *gorm.DB) {
	errors := []error{}

	if err := views.ViewUserWallets(db); err != nil {
		errors = append(errors, err)
	}
	if err := views.ViewUserInvestments(db); err != nil {
		errors = append(errors, err)
	}
	if err := views.ViewUserTransactions(db); err != nil {
		errors = append(errors, err)
	}
	if err := views.ViewUserWalletsGroupByType(db); err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		for _, err := range errors {
			log.Println("Error creating view:", err)
		}
	} else {
		log.Println("All views created successfully")
	}
}
