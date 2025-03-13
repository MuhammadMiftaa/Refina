package repository

import (
	"context"
	"errors"
	
	"server/internal/entity"
	
	"gorm.io/gorm"
)

type TransactionsRepository interface {
	GetAllTransactions(ctx context.Context, tx Transaction) ([]entity.Transactions, error)
	GetTransactionByID(ctx context.Context, tx Transaction, id string) (entity.Transactions, error)
	GetTransactionsByWalletID(ctx context.Context, tx Transaction, id string) ([]entity.Transactions, error)
	CreateTransaction(ctx context.Context, tx Transaction, transaction entity.Transactions) (entity.Transactions, error)
	UpdateTransaction(ctx context.Context, tx Transaction, transaction entity.Transactions) (entity.Transactions, error)
	DeleteTransaction(ctx context.Context, tx Transaction, transaction entity.Transactions) (entity.Transactions, error)
}

type transactionsRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionsRepository {
	return &transactionsRepository{db}
}

func (transaction_repo *transactionsRepository) getDB(ctx context.Context, tx Transaction) (*gorm.DB, error) {
	if tx != nil {
		gormTx, ok := tx.(*GormTx)
		if !ok {
			return nil, errors.New("invalid transaction type")
		}
		return gormTx.db.WithContext(ctx), nil
	}
	return transaction_repo.db.WithContext(ctx), nil
}

func (transaction_repo *transactionsRepository) GetAllTransactions(ctx context.Context, tx Transaction) ([]entity.Transactions, error) {
	db, err := transaction_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var transactions []entity.Transactions
	if err := db.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (transaction_repo *transactionsRepository) GetTransactionByID(ctx context.Context, tx Transaction, id string) (entity.Transactions, error) {
	db, err := transaction_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Transactions{}, err
	}

	var transaction entity.Transactions
	if err := db.First(&transaction, "id = ?", id).Error; err != nil {
		return entity.Transactions{}, errors.New("transaction not found")
	}
	return transaction, nil
}

func (transaction_repo *transactionsRepository) GetTransactionsByWalletID(ctx context.Context, tx Transaction, id string) ([]entity.Transactions, error) {
	db, err := transaction_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var transactions []entity.Transactions
	if err := db.Find(&transactions, "wallet_id = ?", id).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (transaction_repo *transactionsRepository) CreateTransaction(ctx context.Context, tx Transaction, transaction entity.Transactions) (entity.Transactions, error) {
	db, err := transaction_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Transactions{}, err
	}

	if err := db.Create(&transaction).Error; err != nil {
		return entity.Transactions{}, err
	}

	return transaction, nil
}

func (transaction_repo *transactionsRepository) UpdateTransaction(ctx context.Context, tx Transaction, transaction entity.Transactions) (entity.Transactions, error) {
	db, err := transaction_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Transactions{}, err
	}

	if err := db.Save(&transaction).Error; err != nil {
		return entity.Transactions{}, err
	}
	return transaction, nil
}

func (transaction_repo *transactionsRepository) DeleteTransaction(ctx context.Context, tx Transaction, transaction entity.Transactions) (entity.Transactions, error) {
	db, err := transaction_repo.getDB(ctx, tx)
	if err != nil {
		return entity.Transactions{}, err
	}

	if err := db.Delete(&transaction).Error; err != nil {
		return entity.Transactions{}, err
	}
	return transaction, nil
}