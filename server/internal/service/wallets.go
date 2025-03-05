package service

import (
	"errors"

	"server/internal/entity"
	"server/internal/repository"

	"github.com/google/uuid"
)

type WalletsService interface {
	GetAllWallets() ([]entity.Wallets, error)
	GetWalletByID(id string) (entity.Wallets, error)
	GetWalletsByUserID(id string) ([]entity.Wallets, error)
	CreateWallet(wallet entity.WalletsRequest) (entity.Wallets, error)
	UpdateWallet(id string, wallet entity.WalletsRequest) (entity.Wallets, error)
	DeleteWallet(id string) (entity.Wallets, error)
}

type walletsService struct {
	walletsRepository repository.WalletsRepository
}

func NewWalletsService(walletsRepository repository.WalletsRepository) WalletsService {
	return &walletsService{walletsRepository}
}

func (wallet_serv *walletsService) GetAllWallets() ([]entity.Wallets, error) {
	return wallet_serv.walletsRepository.GetAllWallets()
}

func (wallet_serv *walletsService) GetWalletByID(id string) (entity.Wallets, error) {
	return wallet_serv.walletsRepository.GetWalletByID(id)
}

func (wallet_serv *walletsService) GetWalletsByUserID(id string) ([]entity.Wallets, error) {
	return wallet_serv.walletsRepository.GetWalletsByUserID(id)
}

func (wallet_serv *walletsService) CreateWallet(wallet entity.WalletsRequest) (entity.Wallets, error) {
	return wallet_serv.walletsRepository.CreateWallet(entity.Wallets{
		UserID:       uuid.MustParse(wallet.UserID),
		WalletTypeID: uuid.MustParse(wallet.WalletTypeID),
		Name:         wallet.Name,
		Number:       wallet.Number,
		Balance:      wallet.Balance,
	})
}

func (wallet_serv *walletsService) UpdateWallet(id string, wallet entity.WalletsRequest) (entity.Wallets, error) {
	existingWallet, err := wallet_serv.walletsRepository.GetWalletByID(id)
	if err != nil {
		return entity.Wallets{}, errors.New("wallet not found")
	}

	existingWallet.Name = wallet.Name
	existingWallet.Number = wallet.Number
	existingWallet.Balance = wallet.Balance

	return wallet_serv.walletsRepository.UpdateWallet(existingWallet)
}

func (wallet_serv *walletsService) DeleteWallet(id string) (entity.Wallets, error) {
	existingWallet, err := wallet_serv.walletsRepository.GetWalletByID(id)
	if err != nil {
		return entity.Wallets{}, errors.New("wallet not found")
	}

	return wallet_serv.walletsRepository.DeleteWallet(existingWallet)
}