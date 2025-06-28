package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"server/internal/types/entity"
	"server/internal/types/view"

	"gorm.io/gorm"
)

type WalletsRepository interface {
	GetAllWallets(ctx context.Context, tx Transaction) ([]entity.Wallets, error)
	GetWalletByID(ctx context.Context, tx Transaction, id string) (entity.Wallets, error)
	GetWalletsByUserID(ctx context.Context, tx Transaction, id string) ([]view.ViewUserWallets, error)
	GetWalletsByUserIDGroupByType(ctx context.Context, tx Transaction, id string) ([]view.ViewUserWalletsGroupByType, error)
	CreateWallet(ctx context.Context, tx Transaction, wallet entity.Wallets) (entity.Wallets, error)
	UpdateWallet(ctx context.Context, tx Transaction, wallet entity.Wallets) (entity.Wallets, error)
	DeleteWallet(ctx context.Context, tx Transaction, wallet entity.Wallets) (entity.Wallets, error)
}

type walletsRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletsRepository {
	return &walletsRepository{db}
}

// Helper untuk mendapatkan DB instance (transaksi atau biasa)
func (wallet_repo *walletsRepository) getDB(ctx context.Context, tx Transaction) (*gorm.DB, error) {
	if tx != nil {
		gormTx, ok := tx.(*GormTx) // Type assertion ke GORM transaction
		if !ok {
			return nil, errors.New("invalid transaction type")
		}
		return gormTx.db.WithContext(ctx), nil
	}
	return wallet_repo.db.WithContext(ctx), nil
}

// Implementasi method dengan transaksi opsional
func (wallet_repo *walletsRepository) GetAllWallets(ctx context.Context, tx Transaction) ([]entity.Wallets, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var wallets []entity.Wallets
	if err := db.Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

func (wallet_repo *walletsRepository) GetWalletByID(ctx context.Context, tx Transaction, id string) (entity.Wallets, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Wallets{}, err
	}

	var wallet entity.Wallets
	if err := db.Where("id = ?", id).First(&wallet).Error; err != nil {
		return entity.Wallets{}, err
	}
	return wallet, nil
}

func (wallet_repo *walletsRepository) GetWalletsByUserID(ctx context.Context, tx Transaction, id string) ([]view.ViewUserWallets, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var userWallets []view.ViewUserWallets
	err = db.Table("view_user_wallets").Where("user_id = ?", id).Find(&userWallets).Error
	if err != nil {
		return nil, errors.New("user wallets not found")
	}

	if len(userWallets) == 0 {
		return nil, nil
	}

	return userWallets, nil
}

func (wallet_repo *walletsRepository) GetWalletsByUserIDGroupByType(ctx context.Context, tx Transaction, id string) ([]view.ViewUserWalletsGroupByType, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var rawResults []struct {
		UserID     string
		Type string
		Wallets []byte
	}
	err = db.Raw(`SELECT * FROM view_user_wallets_group_by_type WHERE user_id = $1`, id).Scan(&rawResults).Error
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, errors.New("user wallets group by type not found")
	}

	var results []view.ViewUserWalletsGroupByType

	for _, row := range rawResults {
		var wallets []view.ViewUserWalletsGroupByTypeDetailWallet

		err := json.Unmarshal(row.Wallets, &wallets)
		if err != nil {
			return nil, fmt.Errorf("gagal decode JSON wallets (type: %s): %w", row.Type, err)
		}

		results = append(results, view.ViewUserWalletsGroupByType{
			UserID:  row.UserID,
			Type:    row.Type,
			Wallets: wallets,
		})
	}

	return results, nil
}

func (wallet_repo *walletsRepository) CreateWallet(ctx context.Context, tx Transaction, wallet entity.Wallets) (entity.Wallets, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Wallets{}, err
	}

	if err := db.Create(&wallet).Error; err != nil {
		return entity.Wallets{}, err
	}

	return wallet, nil
}

func (wallet_repo *walletsRepository) UpdateWallet(ctx context.Context, tx Transaction, wallet entity.Wallets) (entity.Wallets, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Wallets{}, err
	}

	if err := db.Save(&wallet).Error; err != nil {
		return entity.Wallets{}, err
	}

	return wallet, nil
}

func (wallet_repo *walletsRepository) DeleteWallet(ctx context.Context, tx Transaction, wallet entity.Wallets) (entity.Wallets, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Wallets{}, err
	}

	if err := db.Delete(&wallet).Error; err != nil {
		return entity.Wallets{}, err
	}

	return wallet, nil
}
