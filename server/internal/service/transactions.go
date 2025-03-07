package service

import (
	"errors"

	"server/internal/entity"
	"server/internal/helper"
	"server/internal/repository"
)

type TransactionsService interface {
	GetAllTransactions() ([]entity.Transactions, error)
	GetTransactionByID(id string) (entity.Transactions, error)
	GetTransactionsByWalletID(id string) ([]entity.Transactions, error)
	CreateTransaction(transaction entity.TransactionsRequest) (entity.Transactions, error)
	UpdateTransaction(id string, transaction entity.TransactionsRequest) (entity.Transactions, error)
	DeleteTransaction(id string) (entity.Transactions, error)
}

type transactionsService struct {
	transactionRepo repository.TransactionsRepository
	walletRepo      repository.WalletsRepository
	categoryRepo    repository.CategoriesRepository
}

func NewTransactionService(transactionRepo repository.TransactionsRepository, walletRepo repository.WalletsRepository, categoryRepo repository.CategoriesRepository) TransactionsService {
	return &transactionsService{
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
		categoryRepo:    categoryRepo,
	}
}

func (transaction_serv *transactionsService) GetAllTransactions() ([]entity.Transactions, error) {
	transactions, err := transaction_serv.transactionRepo.GetAllTransactions()
	if err != nil {
		return nil, errors.New("failed to get transactions")
	}

	return transactions, nil
}

func (transaction_serv *transactionsService) GetTransactionByID(id string) (entity.Transactions, error) {
	transaction, err := transaction_serv.transactionRepo.GetTransactionByID(id)
	if err != nil {
		return entity.Transactions{}, errors.New("transaction not found")
	}

	return transaction, nil
}

func (transaction_serv *transactionsService) GetTransactionsByWalletID(id string) ([]entity.Transactions, error) {
	transactions, err := transaction_serv.transactionRepo.GetTransactionsByWalletID(id)
	if err != nil {
		return nil, errors.New("failed to get transactions")
	}

	return transactions, nil
}

func (transaction_serv *transactionsService) CreateTransaction(transaction entity.TransactionsRequest) (entity.Transactions, error) {
	// Check if wallet and category exist
	if _, err := transaction_serv.walletRepo.GetWalletByID(transaction.WalletID); err != nil {
		return entity.Transactions{}, errors.New("wallet not found")
	}

	if _, err := transaction_serv.categoryRepo.GetCategoryByID(transaction.CategoryID); err != nil {
		return entity.Transactions{}, errors.New("category not found")
	}

	// Parse ID from JSON to valid UUID
	CategoryID, err := helper.ParseUUID(transaction.CategoryID)
	if err != nil {
		return entity.Transactions{}, errors.New("invalid category id")
	}

	WalletID, err := helper.ParseUUID(transaction.WalletID)
	if err != nil {
		return entity.Transactions{}, errors.New("invalid wallet id")
	}

	transactionNew, err := transaction_serv.transactionRepo.CreateTransaction(entity.Transactions{
		WalletID:        WalletID,
		CategoryID:      CategoryID,
		Amount:          transaction.Amount,
		TransactionDate: transaction.TransactionDate,
		Description:     transaction.Description,
	})
	if err != nil {
		return entity.Transactions{}, errors.New("failed to create transaction")
	}

	return transactionNew, nil
}

func (transaction_serv *transactionsService) UpdateTransaction(id string, transaction entity.TransactionsRequest) (entity.Transactions, error) {
	// Check if transaction exist
	transactionExist, err := transaction_serv.transactionRepo.GetTransactionByID(id)
	if err != nil {
		return entity.Transactions{}, errors.New("transaction not found")
	}

	// Update transaction only if the field is not empty
	if transaction.CategoryID != "" {
		CategoryID, err := helper.ParseUUID(transaction.CategoryID)
		if err != nil {
			return entity.Transactions{}, errors.New("invalid category id")
		}

		transactionExist.CategoryID = CategoryID
	}
	if transaction.WalletID != "" {
		WalletID, err := helper.ParseUUID(transaction.WalletID)
		if err != nil {
			return entity.Transactions{}, errors.New("invalid wallet id")
		}

		transactionExist.WalletID = WalletID
	}
	if transaction.Amount != 0 {
		transactionExist.Amount = transaction.Amount
	}
	if !transaction.TransactionDate.IsZero() {
		transactionExist.TransactionDate = transaction.TransactionDate
	}
	if transaction.Description != "" {
		transactionExist.Description = transaction.Description
	}

	transactionUpdated, err := transaction_serv.transactionRepo.UpdateTransaction(transactionExist)
	if err != nil {
		return entity.Transactions{}, errors.New("failed to update transaction")
	}

	return transactionUpdated, nil
}

func (transaction_serv *transactionsService) DeleteTransaction(id string) (entity.Transactions, error) {
	// Check if transaction exist
	transactionExist, err := transaction_serv.transactionRepo.GetTransactionByID(id)
	if err != nil {
		return entity.Transactions{}, errors.New("transaction not found")
	}

	transactionDeleted, err := transaction_serv.transactionRepo.DeleteTransaction(transactionExist)
	if err != nil {
		return entity.Transactions{}, errors.New("failed to delete transaction")
	}

	return transactionDeleted, nil
}
