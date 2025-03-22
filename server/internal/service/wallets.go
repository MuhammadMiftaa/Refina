package service

import (
	"context"
	"errors"

	"server/internal/dto"
	"server/internal/entity"
	"server/internal/helper"
	"server/internal/repository"
)

type WalletsService interface {
	GetAllWallets(ctx context.Context) ([]dto.WalletsResponse, error)
	GetWalletByID(ctx context.Context, id string) (dto.WalletsResponse, error)
	GetWalletsByUserID(ctx context.Context, id string) ([]dto.WalletsResponse, error)
	CreateWallet(ctx context.Context, wallet dto.WalletsRequest) (dto.WalletsResponse, error)
	UpdateWallet(ctx context.Context, id string, wallet dto.WalletsRequest) (dto.WalletsResponse, error)
	DeleteWallet(ctx context.Context, id string) (dto.WalletsResponse, error)
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

func (wallet_serv *walletsService) GetAllWallets(ctx context.Context) ([]dto.WalletsResponse, error) {
	wallets, err := wallet_serv.walletsRepository.GetAllWallets(ctx, nil)
	if err != nil {
		return nil, errors.New("failed to get wallets")
	}

	var walletsResponse []dto.WalletsResponse
	for _, wallet := range wallets {
		walletResponse := helper.ConvertToResponseType(wallet).(dto.WalletsResponse)
		walletsResponse = append(walletsResponse, walletResponse)
	}
	
	return walletsResponse, nil
}

func (wallet_serv *walletsService) GetWalletByID(ctx context.Context, id string) (dto.WalletsResponse, error) {
	wallet, err := wallet_serv.walletsRepository.GetWalletByID(ctx, nil, id)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("wallet not found")
	}

	walletResponse := helper.ConvertToResponseType(wallet).(dto.WalletsResponse)

	return walletResponse, nil
}

func (wallet_serv *walletsService) GetWalletsByUserID(ctx context.Context, id string) ([]dto.WalletsResponse, error) {
	wallets, err := wallet_serv.walletsRepository.GetWalletsByUserID(ctx, nil, id)
	if err != nil {
		return nil, errors.New("failed to get wallets")
	}

	var walletsResponse []dto.WalletsResponse
	for _, wallet := range wallets {
		walletResponse := helper.ConvertToResponseType(wallet).(dto.WalletsResponse)
		walletsResponse = append(walletsResponse, walletResponse)
	}
	
	return walletsResponse, err
}

func (wallet_serv *walletsService) CreateWallet(ctx context.Context, wallet dto.WalletsRequest) (dto.WalletsResponse, error) {
	UserID, err := helper.ParseUUID(wallet.UserID)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("invalid user id")
	}

	WalletTypeID, err := helper.ParseUUID(wallet.WalletTypeID)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("invalid wallet type id")
	}

	newWallet, err := wallet_serv.walletsRepository.CreateWallet(ctx, nil, entity.Wallets{
		UserID:       UserID,
		WalletTypeID: WalletTypeID,
		Name:         wallet.Name,
		Number:       wallet.Number,
		Balance:      wallet.Balance,
	})
	if err != nil {
		return dto.WalletsResponse{}, err
	}

	walletResponse := helper.ConvertToResponseType(newWallet).(dto.WalletsResponse)

	return walletResponse, nil
}

func (wallet_serv *walletsService) UpdateWallet(ctx context.Context, id string, wallet dto.WalletsRequest) (dto.WalletsResponse, error) {
	existingWallet, err := wallet_serv.walletsRepository.GetWalletByID(ctx, nil, id)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("wallet not found")
	}

	existingWallet.Name = wallet.Name
	existingWallet.Number = wallet.Number
	existingWallet.Balance = wallet.Balance

	walletUpdated, err := wallet_serv.walletsRepository.UpdateWallet(ctx, nil, existingWallet)
	if err != nil {
		return dto.WalletsResponse{}, err
	}

	walletResponse := helper.ConvertToResponseType(walletUpdated).(dto.WalletsResponse)

	return walletResponse, nil
}

func (wallet_serv *walletsService) DeleteWallet(ctx context.Context, id string) (dto.WalletsResponse, error) {
	existingWallet, err := wallet_serv.walletsRepository.GetWalletByID(ctx, nil, id)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("wallet not found")
	}

	deletedWallet, err := wallet_serv.walletsRepository.DeleteWallet(ctx, nil, existingWallet)
	if err != nil {
		return dto.WalletsResponse{}, err
	}

	walletResponse := helper.ConvertToResponseType(deletedWallet).(dto.WalletsResponse)

	return walletResponse, nil
}
