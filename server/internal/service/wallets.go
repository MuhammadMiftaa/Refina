package service

import (
	"context"
	"errors"

	"server/internal/entity"
	"server/internal/helper"
	"server/internal/repository"
)

type WalletsService interface {
	GetAllWallets(ctx context.Context) ([]entity.Wallets, error)
	GetWalletByID(ctx context.Context, id string) (entity.Wallets, error)
	GetWalletsByUserID(ctx context.Context, id string) ([]entity.Wallets, error)
	CreateWallet(ctx context.Context, wallet entity.WalletsRequest) (entity.Wallets, error)
	UpdateWallet(ctx context.Context, id string, wallet entity.WalletsRequest) (entity.Wallets, error)
	DeleteWallet(ctx context.Context, id string) (entity.Wallets, error)
}

type walletsService struct {
	txManager         repository.TxManager
	walletsRepository repository.WalletsRepository
}

func NewWalletService(txManager repository.TxManager, walletsRepository repository.WalletsRepository) WalletsService {
	return &walletsService{
		txManager:         txManager,
		walletsRepository: walletsRepository,
	}
}

func (wallet_serv *walletsService) GetAllWallets(ctx context.Context) ([]entity.Wallets, error) {
	wallets, err := wallet_serv.walletsRepository.GetAllWallets(ctx, nil)
	if err != nil {
		return nil, errors.New("failed to get wallets")
	}

	return wallets, nil
}

func (wallet_serv *walletsService) GetWalletByID(ctx context.Context, id string) (entity.Wallets, error) {
	wallet, err := wallet_serv.walletsRepository.GetWalletByID(ctx, nil, id)
	if err != nil {
		return entity.Wallets{}, errors.New("wallet not found")
	}

	return wallet, nil
}

func (wallet_serv *walletsService) GetWalletsByUserID(ctx context.Context, id string) ([]entity.Wallets, error) {
	wallets, err := wallet_serv.walletsRepository.GetWalletsByUserID(ctx, nil, id)
	if err != nil {
		return nil, errors.New("failed to get wallets")
	}

	return wallets, err
}

func (wallet_serv *walletsService) CreateWallet(ctx context.Context, wallet entity.WalletsRequest) (entity.Wallets, error) {
	UserID, err := helper.ParseUUID(wallet.UserID)
	if err != nil {
		return entity.Wallets{}, errors.New("invalid user id")
	}

	WalletTypeID, err := helper.ParseUUID(wallet.WalletTypeID)
	if err != nil {
		return entity.Wallets{}, errors.New("invalid wallet type id")
	}

	newWallet, err := wallet_serv.walletsRepository.CreateWallet(ctx, nil, entity.Wallets{
		UserID:       UserID,
		WalletTypeID: WalletTypeID,
		Name:         wallet.Name,
		Number:       wallet.Number,
		Balance:      wallet.Balance,
	})
	if err != nil {
		return entity.Wallets{}, err
	}

	return newWallet, nil
}

func (wallet_serv *walletsService) UpdateWallet(ctx context.Context, id string, wallet entity.WalletsRequest) (entity.Wallets, error) {
	existingWallet, err := wallet_serv.walletsRepository.GetWalletByID(ctx, nil, id)
	if err != nil {
		return entity.Wallets{}, errors.New("wallet not found")
	}

	existingWallet.Name = wallet.Name
	existingWallet.Number = wallet.Number
	existingWallet.Balance = wallet.Balance

	newWallet, err := wallet_serv.walletsRepository.UpdateWallet(ctx, nil, existingWallet)
	if err != nil {
		return entity.Wallets{}, err
	}

	return newWallet, nil
}

func (wallet_serv *walletsService) DeleteWallet(ctx context.Context, id string) (entity.Wallets, error) {
	existingWallet, err := wallet_serv.walletsRepository.GetWalletByID(ctx, nil, id)
	if err != nil {
		return entity.Wallets{}, errors.New("wallet not found")
	}

	deletedWallet, err := wallet_serv.walletsRepository.DeleteWallet(ctx, nil, existingWallet)
	if err != nil {
		return entity.Wallets{}, err
	}

	return deletedWallet, nil
}
