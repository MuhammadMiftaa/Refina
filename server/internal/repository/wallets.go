package repository

import (
	"context"
	"errors"

	"server/internal/entity"

	"gorm.io/gorm"
)

type WalletsRepository interface {
	GetAllWallets(ctx context.Context, tx Transaction) ([]entity.Wallets, error)
	GetWalletByID(ctx context.Context, tx Transaction, id string) (entity.Wallets, error)
	GetWalletsByUserID(ctx context.Context, tx Transaction, id string) ([]entity.Wallets, error)
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

func (wallet_repo *walletsRepository) GetWalletsByUserID(ctx context.Context, tx Transaction, id string) ([]entity.Wallets, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var wallets []entity.Wallets
	if err := db.Where("user_id = ?", id).Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
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