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
	GetTransactionsByUserID(ctx context.Context, tx Transaction, id string) ([]entity.TransactionsList, error)
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

func (transaction_repo *transactionsRepository) GetTransactionsByUserID(ctx context.Context, tx Transaction, id string) ([]entity.TransactionsList, error) {
	db, err := transaction_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var transactions []entity.TransactionsList
	err = db.Table("users").Select("wallets.id AS wallet_id, wallets.name AS wallet_name, wallets.number AS wallet_number, wallets.balance AS wallet_balance, wallet_types.name AS wallet_type, transactions.id AS transaction_id, categories.name AS category_name, categories.type AS category_type, transactions.amount, transactions.transaction_date, transactions.description, attachments.image").
		Joins("LEFT JOIN wallets ON users.id = wallets.user_id AND wallets.deleted_at IS NULL").
		Joins("LEFT JOIN wallet_types ON wallets.wallet_type_id = wallet_types.id AND wallet_types.deleted_at IS NULL").
		Joins("LEFT JOIN transactions ON wallets.id = transactions.wallet_id AND transactions.deleted_at IS NULL").
		Joins("LEFT JOIN categories ON transactions.category_id = categories.id AND categories.deleted_at IS NULL").
		Joins("LEFT JOIN attachments ON transactions.id = attachments.transaction_id AND attachments.deleted_at IS NULL").
		Where("users.id = ?", id).
		Where("users.deleted_at IS NULL").
		Find(&transactions).Error
	if err != nil {
		return nil, errors.New("user transactions not found")
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
