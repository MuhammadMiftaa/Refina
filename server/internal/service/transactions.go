package service

import (
	"errors"

	"github.com/google/uuid"

	"server/internal/entity"
	"server/internal/repository"
)

type TransactionsService interface {
	GetAllTransactions() ([]entity.Transactions, error)
	GetTransactionByID(id string) (entity.Transactions, error)
	CreateTransaction(transaction entity.TransactionsRequest) (entity.Transactions, error)
	UpdateTransaction(id string, transaction entity.TransactionsRequest) (entity.Transactions, error)
	DeleteTransaction(id string) (entity.Transactions, error)
}

type transactionsService struct {
	transactionRepo repository.TransactionsRepository
	userRepo        repository.UsersRepository
}

func NewTransactionService(transactionRepo repository.TransactionsRepository, userRepo repository.UsersRepository) TransactionsService {
	return &transactionsService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

func (transactionServ *transactionsService) GetAllTransactions() ([]entity.Transactions, error) {
	transactions, err := transactionServ.transactionRepo.GetAllTransactions()
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (transactionServ *transactionsService) GetTransactionByID(id string) (entity.Transactions, error) {
	transaction, err := transactionServ.transactionRepo.GetTransactionByID(id)
	if err != nil {
		return entity.Transactions{}, err
	}

	return transaction, nil
}

func (transactionServ *transactionsService) CreateTransaction(transaction entity.TransactionsRequest) (entity.Transactions, error) {
	if transaction.Amount == 0 || transaction.CategoryID == "" || transaction.Description == "" || transaction.WalletID == "" {
		return entity.Transactions{}, errors.New("all fields must be filled")
	}

	if transaction.Amount < 0 {
		return entity.Transactions{}, errors.New("amount must be greater than 0")
	}

	// if transaction.CategoryID != "income" && transaction.TransactionType != "expense" {
	// 	return entity.Transactions{}, errors.New("transaction type must be income or expense")
	// }

	if _, err := transactionServ.userRepo.GetUserByID(transaction.WalletID); err != nil {
		return entity.Transactions{}, errors.New("user not found")
	}

	transactionNew := entity.Transactions{
		WalletID:        uuid.MustParse(transaction.WalletID),
		CategoryID:      uuid.MustParse(transaction.CategoryID),
		Amount:          transaction.Amount,
		TransactionDate: transaction.TransactionDate,
		Description:     transaction.Description,
	}

	transactionCreated, err := transactionServ.transactionRepo.CreateTransaction(transactionNew)
	if err != nil {
		return entity.Transactions{}, err
	}

	return transactionCreated, nil
}

func (transactionServ *transactionsService) UpdateTransaction(id string, transaction entity.TransactionsRequest) (entity.Transactions, error) {
	transactionExist, err := transactionServ.transactionRepo.GetTransactionByID(id)
	if err != nil {
		return entity.Transactions{}, errors.New("transaction not found")
	}

	transactionExist.Amount = transaction.Amount
	transactionExist.CategoryID = uuid.MustParse(transaction.CategoryID)
	transactionExist.Description = transaction.Description
	transactionExist.WalletID = uuid.MustParse(transaction.WalletID)
	transactionExist.TransactionDate = transaction.TransactionDate

	transactionUpdated, err := transactionServ.transactionRepo.UpdateTransaction(transactionExist)
	if err != nil {
		return entity.Transactions{}, err
	}

	return transactionUpdated, nil
}

func (transactionServ *transactionsService) DeleteTransaction(id string) (entity.Transactions, error) {
	transactionExist, err := transactionServ.transactionRepo.GetTransactionByID(id)
	if err != nil {
		return entity.Transactions{}, errors.New("transaction not found")
	}

	transactionDeleted, err := transactionServ.transactionRepo.DeleteTransaction(transactionExist)
	if err != nil {
		return entity.Transactions{}, err
	}

	return transactionDeleted, nil
}
