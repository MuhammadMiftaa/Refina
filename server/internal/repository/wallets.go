package repository

import (
	"server/internal/entity"

	"gorm.io/gorm"
)

type WalletsRepository interface {
	GetAllWallets() ([]entity.Wallets, error)
	GetWalletByID(id string) (entity.Wallets, error)
	GetWalletsByUserID(id string) ([]entity.Wallets, error)
	CreateWallet(wallet entity.Wallets) (entity.Wallets, error)
	UpdateWallet(wallet entity.Wallets) (entity.Wallets, error)
	DeleteWallet(wallet entity.Wallets) (entity.Wallets, error)
}

type walletsRepository struct {
	db *gorm.DB
}

func NewWalletsRepository(db *gorm.DB) WalletsRepository {
	return &walletsRepository{db}
}

func (wallet_repo *walletsRepository) GetAllWallets() ([]entity.Wallets, error) {
	var wallets []entity.Wallets
	err := wallet_repo.db.Find(&wallets).Error
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func (wallet_repo *walletsRepository) GetWalletByID(id string) (entity.Wallets, error) {
	var wallet entity.Wallets
	err := wallet_repo.db.First(&wallet, "id = ?", id).Error
	if err != nil {
		return entity.Wallets{}, err
	}

	return wallet, nil
}

func (wallet_repo *walletsRepository) GetWalletsByUserID(id string) ([]entity.Wallets, error) {
	var wallets []entity.Wallets
	err := wallet_repo.db.Find(&wallets, "user_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func (wallet_repo *walletsRepository) CreateWallet(wallet entity.Wallets) (entity.Wallets, error) {
	err := wallet_repo.db.Create(&wallet).Error
	if err != nil {
		return entity.Wallets{}, err
	}

	return wallet, nil
}

func (wallet_repo *walletsRepository) UpdateWallet(wallet entity.Wallets) (entity.Wallets, error) {
	err := wallet_repo.db.Save(&wallet).Error
	if err != nil {
		return entity.Wallets{}, err
	}

	return wallet, nil
}

func (wallet_repo *walletsRepository) DeleteWallet(wallet entity.Wallets) (entity.Wallets, error) {
	err := wallet_repo.db.Delete(&wallet).Error
	if err != nil {
		return entity.Wallets{}, err
	}

	return wallet, nil
}
