package service

import (
	"context"
	"errors"
	"testing"

	"server/internal/dto"
	"server/internal/entity"
	"server/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllWallets(t *testing.T) {
	wallet_repo_mock := repository.NewWalletsRepositoryMock()
	wallet_serv_test := NewWalletService(nil, wallet_repo_mock)

	defer wallet_repo_mock.Mock.AssertExpectations(t)

	wallets := []entity.Wallets{
		{
			Base: entity.Base{
				ID: uuid.MustParse("c1b9b26d-3390-4ff9-ae34-838050b52f90"),
			},
			Name:         "Wallet 1",
			Balance:      1000,
			UserID:       uuid.MustParse("c1b9b26d-3390-4ff9-ae34-838050b52f90"),
			WalletTypeID: uuid.MustParse("dd54a719-4f28-43ca-8f3a-998cc62f09ef"),
			Number:       "1234567890",
		},
		{
			Base: entity.Base{
				ID: uuid.MustParse("c1b9b26d-3390-4ff9-ae34-838050b52f90"),
			},
			Name:         "Wallet 2",
			Balance:      2000,
			UserID:       uuid.MustParse("c1b9b26d-3390-4ff9-ae34-838050b52f90"),
			WalletTypeID: uuid.MustParse("dd54a719-4f28-43ca-8f3a-998cc62f09ef"),
			Number:       "0987654321",
		},
	}
	wallet_repo_mock.Mock.On("GetAllWallets", context.Background(), mock.Anything).Return(wallets, nil)
	wallet_repo_mock.Mock.On("GetAllWallets", nil, mock.Anything).Return(nil, errors.New("error getting all wallets, there is context"))

	t.Run("Get All Wallets Success", func(t *testing.T) {
		wallets_response, err := wallet_serv_test.GetAllWallets(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, wallets_response)
		assert.Equal(t, len(wallets_response), 2)
		assert.Equal(t, wallets_response[0].ID, wallets[0].ID.String())
		assert.Equal(t, wallets_response[0].Name, wallets[0].Name)
		assert.Equal(t, wallets_response[0].Balance, wallets[0].Balance)
		assert.Equal(t, wallets_response[0].UserID, wallets[0].UserID.String())
		assert.Equal(t, wallets_response[0].WalletTypeID, wallets[0].WalletTypeID.String())
		assert.Equal(t, wallets_response[0].Number, wallets[0].Number)
		assert.Equal(t, wallets_response[1].ID, wallets[1].ID.String())
		assert.Equal(t, wallets_response[1].Name, wallets[1].Name)
		assert.Equal(t, wallets_response[1].Balance, wallets[1].Balance)
		assert.Equal(t, wallets_response[1].UserID, wallets[1].UserID.String())
		assert.Equal(t, wallets_response[1].WalletTypeID, wallets[1].WalletTypeID.String())
		assert.Equal(t, wallets_response[1].Number, wallets[1].Number)
	})

	t.Run("Get All Wallets Failed", func(t *testing.T) {
		wallets_response, err := wallet_serv_test.GetAllWallets(nil)
		assert.NotNil(t, err)
		assert.Nil(t, wallets_response)
		assert.Equal(t, err.Error(), "failed to get wallets")
	})
}

func TestGetWalletByID(t *testing.T) {
	wallet_repo_mock := repository.NewWalletsRepositoryMock()
	wallet_serv_test := NewWalletService(nil, wallet_repo_mock)

	defer wallet_repo_mock.Mock.AssertExpectations(t)

	wallet := entity.Wallets{
		Base: entity.Base{
			ID: uuid.MustParse("c1b9b26d-3390-4ff9-ae34-838050b52f90"),
		},
		Name:         "Wallet 1",
		Balance:      1000,
		UserID:       uuid.MustParse("c1b9b26d-3390-4ff9-ae34-838050b52f90"),
		WalletTypeID: uuid.MustParse("dd54a719-4f28-43ca-8f3a-998cc62f09ef"),
		Number:       "1234567890",
	}
	wallet_repo_mock.Mock.On("GetWalletByID", context.Background(), mock.Anything, wallet.ID.String()).Return(wallet, nil)
	wallet_repo_mock.Mock.On("GetWalletByID", context.Background(), nil, "123").Return(entity.Wallets{}, errors.New("error getting wallet by ID"))

	t.Run("Get Wallet By ID Success", func(t *testing.T) {
		wallet_response, err := wallet_serv_test.GetWalletByID(context.Background(), wallet.ID.String())
		assert.Nil(t, err)
		assert.NotNil(t, wallet_response)
		assert.Equal(t, wallet_response.ID, wallet.ID.String())
		assert.Equal(t, wallet_response.Name, wallet.Name)
		assert.Equal(t, wallet_response.Balance, wallet.Balance)
		assert.Equal(t, wallet_response.UserID, wallet.UserID.String())
		assert.Equal(t, wallet_response.WalletTypeID, wallet.WalletTypeID.String())
		assert.Equal(t, wallet_response.Number, wallet.Number)
	})

	t.Run("Get Wallet By ID Failed", func(t *testing.T) {
		wallet_response, err := wallet_serv_test.GetWalletByID(context.Background(), "123")
		assert.NotNil(t, err)
		assert.Equal(t, wallet_response, dto.WalletsResponse{})
	})
}
