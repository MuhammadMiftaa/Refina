package repository

import (
	"errors"

	"server/internal/entity"

	"gorm.io/gorm"
)

type TransactionsRepository interface {
	GetAllTransactions() ([]entity.Transactions, error)
	GetTransactionByID(id string) (entity.Transactions, error)
	GetTransactionsByWalletID(id string) ([]entity.Transactions, error)
	CreateTransaction(transaction entity.Transactions) (entity.Transactions, error)
	UpdateTransaction(transaction entity.Transactions) (entity.Transactions, error)
	DeleteTransaction(transaction entity.Transactions) (entity.Transactions, error)
}

type transactionsRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionsRepository {
	return &transactionsRepository{db}
}

func (transaction_repo *transactionsRepository) GetAllTransactions() ([]entity.Transactions, error) {
	var transactions []entity.Transactions
	err := transaction_repo.db.Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (transaction_repo *transactionsRepository) GetTransactionByID(id string) (entity.Transactions, error) {
	var transaction entity.Transactions
	err := transaction_repo.db.First(&transaction, "id = ?", id).Error
	if err != nil {
		return entity.Transactions{}, errors.New("transaction not found")
	}

	return transaction, nil
}

func (transaction_repo *transactionsRepository) GetTransactionsByWalletID(id string) ([]entity.Transactions, error) {
	var transactions []entity.Transactions
	err := transaction_repo.db.Find(&transactions, "wallet_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (transaction_repo *transactionsRepository) CreateTransaction(transaction entity.Transactions) (entity.Transactions, error) {
	err := transaction_repo.db.Create(&transaction).Error
	if err != nil {
		return entity.Transactions{}, err
	}
	return transaction, nil
}

func (transaction_repo *transactionsRepository) UpdateTransaction(transaction entity.Transactions) (entity.Transactions, error) {
	err := transaction_repo.db.Save(&transaction).Error
	if err != nil {
		return entity.Transactions{}, err
	}
	return transaction, nil
}

func (transaction_repo *transactionsRepository) DeleteTransaction(transaction entity.Transactions) (entity.Transactions, error) {
	err := transaction_repo.db.Delete(&transaction).Error
	if err != nil {
		return entity.Transactions{}, err
	}
	return transaction, nil
}
