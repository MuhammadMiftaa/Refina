package config

import (
	"fmt"
	"log"

	"server/db/prepopulate"
	"server/db/views"
	"server/env/config"
	"server/internal/types/entity"

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

	if err := db.AutoMigrate(&entity.Users{}, &entity.WalletTypes{}, &entity.Wallets{}, &entity.Categories{}, &entity.Transactions{}, &entity.Attachments{}, &entity.InvestmentTypes{}, &entity.Investments{}); err != nil {
		log.Fatalf("Error saat melakukan migrasi: %v", err)
		panic(err)
	}

	DB = db

	PrePopulate()
	if err := CreateView(); err != nil {
		log.Fatalf("Error saat membuat view: %v", err)
		panic(err)
	}
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

func PrePopulate() {
	prepopulate.WalletTypesSeeder(DB)
	prepopulate.CategoryTypeSeeder(DB)
	prepopulate.InvestmentTypesSeeder(DB)
}

func CreateView() error {
	errors := []error{}

	if err := views.ViewUserWallets(DB); err != nil {
		errors = append(errors, err)
	}
	if err := views.ViewUserInvestments(DB); err != nil {
		errors = append(errors, err)
	}
	if err := views.ViewUserTransactions(DB); err != nil {
		errors = append(errors, err)
	}
	if err := views.ViewUserWalletsGroupByType(DB); err != nil {
		errors = append(errors, err)
	}
	if err := views.ViewCategoryGroupByType(DB); err != nil {
		errors = append(errors, err)
	}
	if err := views.MVUserSummaries(DB); err != nil {
		errors = append(errors, err)
	}
	if err := views.MVUserMonthlySummary(DB); err != nil {
		errors = append(errors, err)
	}
	if err := views.MVUserMostExpense(DB); err != nil {
		errors = append(errors, err)
	}
	if err := views.MVUserWalletDailySummary(DB); err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		for _, err := range errors {
			log.Println("Error creating view:", err)
		}
		return fmt.Errorf("failed to create views: %v", errors)
	}
	
	return nil
}
