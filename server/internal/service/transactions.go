package service

import (
	"context"
	"errors"

	"server/internal/entity"
	"server/internal/helper"
	"server/internal/repository"
)

type TransactionsService interface {
	GetAllTransactions(ctx context.Context) ([]entity.Transactions, error)
	GetTransactionByID(ctx context.Context, id string) (entity.Transactions, error)
	GetTransactionsByWalletID(ctx context.Context, id string) ([]entity.Transactions, error)
	CreateTransaction(ctx context.Context, transaction entity.TransactionsRequest) (entity.Transactions, error)
	UpdateTransaction(ctx context.Context, id string, transaction entity.TransactionsRequest) (entity.Transactions, error)
	DeleteTransaction(ctx context.Context, id string) (entity.Transactions, error)
}

type transactionsService struct {
	txManager       repository.TxManager
	transactionRepo repository.TransactionsRepository
	walletRepo      repository.WalletsRepository
	categoryRepo    repository.CategoriesRepository
}

func NewTransactionService(txManager repository.TxManager, transactionRepo repository.TransactionsRepository, walletRepo repository.WalletsRepository, categoryRepo repository.CategoriesRepository) TransactionsService {
	return &transactionsService{
		txManager:       txManager,
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
		categoryRepo:    categoryRepo,
	}
}

func (transaction_serv *transactionsService) GetAllTransactions(ctx context.Context) ([]entity.Transactions, error) {
	transactions, err := transaction_serv.transactionRepo.GetAllTransactions(ctx, nil)
	if err != nil {
		return nil, errors.New("failed to get transactions")
	}

	return transactions, nil
}

func (transaction_serv *transactionsService) GetTransactionByID(ctx context.Context, id string) (entity.Transactions, error) {
	transaction, err := transaction_serv.transactionRepo.GetTransactionByID(ctx, nil, id)
	if err != nil {
		return entity.Transactions{}, errors.New("transaction not found")
	}

	return transaction, nil
}

func (transaction_serv *transactionsService) GetTransactionsByWalletID(ctx context.Context, id string) ([]entity.Transactions, error) {
	transactions, err := transaction_serv.transactionRepo.GetTransactionsByWalletID(ctx, nil, id)
	if err != nil {
		return nil, errors.New("failed to get transactions")
	}

	return transactions, nil
}

func (transaction_serv *transactionsService) CreateTransaction(ctx context.Context, transaction entity.TransactionsRequest) (entity.Transactions, error) {
	tx, err := transaction_serv.txManager.Begin(ctx)
	if err != nil {
		return entity.Transactions{}, errors.New("failed to create transaction")
	}

	defer func() {
		// Rollback otomatis jika transaksi belum di-commit
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	// Check if wallet and category exist
	wallet, err := transaction_serv.walletRepo.GetWalletByID(ctx, tx, transaction.WalletID)
	if err != nil {
		return entity.Transactions{}, errors.New("wallet not found")
	}

	category, err := transaction_serv.categoryRepo.GetCategoryByID(ctx, tx, transaction.CategoryID)
	if err != nil {
		return entity.Transactions{}, errors.New("category not found")
	}

	// Check if transaction type is valid and update wallet balance
	if category.Type == "expense" {
		wallet.Balance -= transaction.Amount
	} else if category.Type == "income" {
		wallet.Balance += transaction.Amount
	} else {
		return entity.Transactions{}, errors.New("invalid transaction type")
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

	// Update wallet balance
	_, err = transaction_serv.walletRepo.UpdateWallet(ctx, tx, wallet)
	if err != nil {
		return entity.Transactions{}, errors.New("failed to update wallet")
	}

	// Create transaction
	transactionNew, err := transaction_serv.transactionRepo.CreateTransaction(ctx, tx, entity.Transactions{
		WalletID:        WalletID,
		CategoryID:      CategoryID,
		Amount:          transaction.Amount,
		TransactionDate: transaction.TransactionDate,
		Description:     transaction.Description,
	})
	if err != nil {
		return entity.Transactions{}, errors.New("failed to create transaction")
	}

	// Commit transaksi jika semua sukses
    if err := tx.Commit(); err != nil {
        return entity.Transactions{}, errors.New("failed to commit transaction")
    }

	return transactionNew, nil
}

func (transaction_serv *transactionsService) UpdateTransaction(ctx context.Context, id string, transaction entity.TransactionsRequest) (entity.Transactions, error) {
	// Check if transaction exist
	transactionExist, err := transaction_serv.transactionRepo.GetTransactionByID(ctx, nil, id)
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

	transactionUpdated, err := transaction_serv.transactionRepo.UpdateTransaction(ctx, nil, transactionExist)
	if err != nil {
		return entity.Transactions{}, errors.New("failed to update transaction")
	}

	return transactionUpdated, nil
}

func (transaction_serv *transactionsService) DeleteTransaction(ctx context.Context, id string) (entity.Transactions, error) {
	// Check if transaction exist
	transactionExist, err := transaction_serv.transactionRepo.GetTransactionByID(ctx, nil, id)
	if err != nil {
		return entity.Transactions{}, errors.New("transaction not found")
	}

	transactionDeleted, err := transaction_serv.transactionRepo.DeleteTransaction(ctx, nil, transactionExist)
	if err != nil {
		return entity.Transactions{}, errors.New("failed to delete transaction")
	}

	return transactionDeleted, nil
}
