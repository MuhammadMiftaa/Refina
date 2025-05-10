package repository

import (
	"context"
	"errors"

	"server/internal/entity"

	"github.com/stretchr/testify/mock"
)

type walletsRepositoryMock struct {
	Mock mock.Mock
}

func NewWalletsRepositoryMock() *walletsRepositoryMock {
	return &walletsRepositoryMock{Mock: mock.Mock{}}
}

func (wallet_repo *walletsRepositoryMock) GetAllWallets(ctx context.Context, tx Transaction) ([]entity.Wallets, error) {
	arguments := wallet_repo.Mock.Called(ctx, tx)
	
	result, ok := arguments.Get(0).([]entity.Wallets)
	if !ok {
		return nil, errors.New("error getting all wallets")
	}

	return result, nil
}

func (wallet_repo *walletsRepositoryMock) GetWalletByID(ctx context.Context, tx Transaction, id string) (entity.Wallets, error) {
	arguments := wallet_repo.Mock.Called(ctx, tx, id)
	
	result, ok := arguments.Get(0).(entity.Wallets)
	if !ok || arguments.Get(1) != nil {
		return entity.Wallets{}, errors.New("error getting wallet by ID")
	}

	return result, nil
}

func (wallet_repo *walletsRepositoryMock) GetWalletsByUserID(ctx context.Context, tx Transaction, id string) ([]entity.Wallets, error) {
	arguments := wallet_repo.Mock.Called(ctx, tx, id)
	
	result, ok := arguments.Get(0).([]entity.Wallets)
	if !ok {
		return nil, errors.New("error getting wallets by user ID")
	}

	return result, nil
}

func (wallet_repo *walletsRepositoryMock) CreateWallet(ctx context.Context, tx Transaction, wallet entity.Wallets) (entity.Wallets, error) {
	arguments := wallet_repo.Mock.Called(ctx, tx, wallet)
	
	result, ok := arguments.Get(0).(entity.Wallets)
	if !ok {
		return entity.Wallets{}, errors.New("error creating wallet")
	}

	return result, nil
}

func (wallet_repo *walletsRepositoryMock) UpdateWallet(ctx context.Context, tx Transaction, wallet entity.Wallets) (entity.Wallets, error) {
	arguments := wallet_repo.Mock.Called(ctx, tx, wallet)
	
	result, ok := arguments.Get(0).(entity.Wallets)
	if !ok {
		return entity.Wallets{}, errors.New("error updating wallet")
	}

	return result, nil
}

func (wallet_repo *walletsRepositoryMock) DeleteWallet(ctx context.Context, tx Transaction, wallet entity.Wallets) (entity.Wallets, error) {
	arguments := wallet_repo.Mock.Called(ctx, tx, wallet)
	
	result, ok := arguments.Get(0).(entity.Wallets)
	if !ok {
		return entity.Wallets{}, errors.New("error deleting wallet")
	}

	return result, nil
}
