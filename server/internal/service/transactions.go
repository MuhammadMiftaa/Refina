package service

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"server/internal/dto"
	"server/internal/entity"
	"server/internal/helper"
	"server/internal/repository"

	"github.com/google/uuid"
)

type TransactionsService interface {
	GetAllTransactions(ctx context.Context) ([]dto.TransactionsResponse, error)
	GetTransactionByID(ctx context.Context, id string) (dto.TransactionsResponse, error)
	GetTransactionsByUserID(ctx context.Context, token string) ([]dto.TransactionsResponse, error)
	CreateTransaction(ctx context.Context, transaction dto.TransactionsRequest) (dto.TransactionsResponse, error)
	FundTransfer(ctx context.Context, transaction dto.FundTransferRequest) (dto.FundTransferResponse, error)
	UploadAttachment(ctx context.Context, transactionID string, file multipart.File, handler *multipart.FileHeader) (dto.AttachmentsResponse, error)
	UpdateTransaction(ctx context.Context, id string, transaction dto.TransactionsRequest) (dto.TransactionsResponse, error)
	DeleteTransaction(ctx context.Context, id string) (dto.TransactionsResponse, error)
}

type transactionsService struct {
	txManager       repository.TxManager
	transactionRepo repository.TransactionsRepository
	walletRepo      repository.WalletsRepository
	categoryRepo    repository.CategoriesRepository
	attachmentRepo  repository.AttachmentsRepository
}

func NewTransactionService(txManager repository.TxManager, transactionRepo repository.TransactionsRepository, walletRepo repository.WalletsRepository, categoryRepo repository.CategoriesRepository, attachmentRepo repository.AttachmentsRepository) TransactionsService {
	return &transactionsService{
		txManager:       txManager,
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
		categoryRepo:    categoryRepo,
		attachmentRepo:  attachmentRepo,
	}
}

func (transaction_serv *transactionsService) GetAllTransactions(ctx context.Context) ([]dto.TransactionsResponse, error) {
	transactions, err := transaction_serv.transactionRepo.GetAllTransactions(ctx, nil)
	if err != nil {
		return nil, errors.New("failed to get transactions")
	}

	var transactionsResponse []dto.TransactionsResponse
	for _, transaction := range transactions {
		transactionResponse := helper.ConvertToResponseType(transaction).(dto.TransactionsResponse)
		transactionsResponse = append(transactionsResponse, transactionResponse)
	}

	return transactionsResponse, nil
}

func (transaction_serv *transactionsService) GetTransactionByID(ctx context.Context, id string) (dto.TransactionsResponse, error) {
	transaction, err := transaction_serv.transactionRepo.GetTransactionByIDJoin(ctx, nil, id)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("transaction not found")
	}

	transactionResponse := helper.ConvertToResponseType(transaction).(dto.TransactionsResponse)

	return transactionResponse, nil
}

func (transaction_serv *transactionsService) GetTransactionsByUserID(ctx context.Context, token string) ([]dto.TransactionsResponse, error) {
	userData, err := helper.VerifyToken(token[7:])
	if err != nil {
		return nil, errors.New("invalid token")
	}

	transactions, err := transaction_serv.transactionRepo.GetTransactionsByUserID(ctx, nil, userData.ID)
	if err != nil {
		return nil, errors.New("failed to get transactions")
	}

	var transactionsResponse []dto.TransactionsResponse
	for _, transaction := range transactions {
		transactionResponse := helper.ConvertToResponseType(transaction).(dto.TransactionsResponse)
		transactionsResponse = append(transactionsResponse, transactionResponse)
	}

	return transactionsResponse, nil
}

func (transaction_serv *transactionsService) CreateTransaction(ctx context.Context, transaction dto.TransactionsRequest) (dto.TransactionsResponse, error) {
	tx, err := transaction_serv.txManager.Begin(ctx)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("failed to create transaction")
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
		return dto.TransactionsResponse{}, errors.New("wallet not found")
	}

	category, err := transaction_serv.categoryRepo.GetCategoryByID(ctx, tx, transaction.CategoryID)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("category not found")
	}

	// Check if transaction type is valid and update wallet balance
	switch category.Type {
	case "expense":
		wallet.Balance -= transaction.Amount
	case "income":
		wallet.Balance += transaction.Amount
	default:
		return dto.TransactionsResponse{}, errors.New("invalid transaction type")
	}

	// Parse ID from JSON to valid UUID
	CategoryID, err := helper.ParseUUID(transaction.CategoryID)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("invalid category id")
	}

	WalletID, err := helper.ParseUUID(transaction.WalletID)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("invalid wallet id")
	}

	// Update wallet balance
	_, err = transaction_serv.walletRepo.UpdateWallet(ctx, tx, wallet)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("failed to update wallet")
	}

	// Create transaction
	transactionNew, err := transaction_serv.transactionRepo.CreateTransaction(ctx, tx, entity.Transactions{
		WalletID:        WalletID,
		CategoryID:      CategoryID,
		Amount:          transaction.Amount,
		TransactionDate: transaction.Date,
		Description:     transaction.Description,
	})
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("failed to create transaction")
	}

	// Commit transaksi jika semua sukses
	if err := tx.Commit(); err != nil {
		return dto.TransactionsResponse{}, errors.New("failed to commit transaction")
	}

	transactionResponse := helper.ConvertToResponseType(transactionNew).(dto.TransactionsResponse)

	return transactionResponse, nil
}

func (transaction_serv *transactionsService) FundTransfer(ctx context.Context, transaction dto.FundTransferRequest) (dto.FundTransferResponse, error) {
	tx, err := transaction_serv.txManager.Begin(ctx)
	if err != nil {
		return dto.FundTransferResponse{}, errors.New("failed to create transaction")
	}

	defer func() {
		// Rollback otomatis jika transaksi belum di-commit
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	// Check if wallet and category exist
	fromWallet, err := transaction_serv.walletRepo.GetWalletByID(ctx, tx, transaction.FromWalletID)
	if err != nil {
		return dto.FundTransferResponse{}, errors.New("source wallet not found")
	}

	toWallet, err := transaction_serv.walletRepo.GetWalletByID(ctx, tx, transaction.ToWalletID)
	if err != nil {
		return dto.FundTransferResponse{}, errors.New("destination wallet not found")
	}

	if fromWallet.ID == toWallet.ID {
		return dto.FundTransferResponse{}, errors.New("source wallet and destination wallet cannot be the same")
	}

	fromWallet.Balance -= (transaction.Amount + transaction.AdminFee)
	toWallet.Balance += transaction.Amount

	// Parse ID from JSON to valid UUID
	FromWalletID, err := helper.ParseUUID(transaction.FromWalletID)
	if err != nil {
		return dto.FundTransferResponse{}, errors.New("invalid from wallet id")
	}

	ToWalletID, err := helper.ParseUUID(transaction.ToWalletID)
	if err != nil {
		return dto.FundTransferResponse{}, errors.New("invalid to wallet id")
	}

	// Parse CategoryID from JSON to valid UUID
	FromCategoryID, err := helper.ParseUUID(transaction.CashOutCategoryID)
	if err != nil {
		return dto.FundTransferResponse{}, errors.New("invalid from category id")
	}

	ToCategoryID, err := helper.ParseUUID(transaction.CashInCategoryID)
	if err != nil {
		return dto.FundTransferResponse{}, errors.New("invalid to category id")
	}

	// Update wallet balance
	if _, err = transaction_serv.walletRepo.UpdateWallet(ctx, tx, fromWallet); err != nil {
		return dto.FundTransferResponse{}, errors.New("failed to update from wallet")
	}
	if _, err = transaction_serv.walletRepo.UpdateWallet(ctx, tx, toWallet); err != nil {
		return dto.FundTransferResponse{}, errors.New("failed to update to wallet")
	}

	transactionNewFrom, err := transaction_serv.transactionRepo.CreateTransaction(ctx, tx, entity.Transactions{
		WalletID:        FromWalletID,
		CategoryID:      FromCategoryID,
		Amount:          transaction.Amount + transaction.AdminFee,
		TransactionDate: transaction.Date,
		Description:     "fund transfer to " + toWallet.Name + "(Cash Out)",
	})
	if err != nil {
		return dto.FundTransferResponse{}, errors.New("failed to create from transaction")
	}

	transactionNewTo, err := transaction_serv.transactionRepo.CreateTransaction(ctx, tx, entity.Transactions{
		WalletID:        ToWalletID,
		CategoryID:      ToCategoryID,
		Amount:          transaction.Amount,
		TransactionDate: transaction.Date,
		Description:     "fund transfer from " + fromWallet.Name + "(Cash In)",
	})
	if err != nil {
		return dto.FundTransferResponse{}, errors.New("failed to create to transaction")
	}

	if err := tx.Commit(); err != nil {
		return dto.FundTransferResponse{}, errors.New("failed to commit transaction")
	}

	response := dto.FundTransferResponse{
		CashOutTransactionID: transactionNewFrom.ID.String(),
		CashInTransactionID:  transactionNewTo.ID.String(),
		FromWalletID:         transaction.FromWalletID,
		ToWalletID:           transaction.ToWalletID,
		Amount:               transaction.Amount,
		Date:                 transaction.Date,
		Description:          transaction.Description,
	}

	return response, nil
}

func (transaction_serv *transactionsService) UploadAttachment(ctx context.Context, transactionID string, file multipart.File, handler *multipart.FileHeader) (dto.AttachmentsResponse, error) {
	if handler.Size > int64(helper.ATTACHMENT_MAX_SIZE) {
		return dto.AttachmentsResponse{}, errors.New("file size exceeds 10MB")
	}

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if !helper.ATTACHMENT_EXT_ALLOWED[ext] {
		return dto.AttachmentsResponse{}, errors.New("invalid file type")
	}

	absolutePath, _ := helper.ExpandPathAndCreateDir(helper.ATTACHMENT_FILEPATH)

	if err := helper.StorageIsExist(absolutePath); err != nil {
		return dto.AttachmentsResponse{}, errors.New("storage not found")
	}
	// Pastikan file bisa dibaca ulang
	file.Seek(0, 0)

	// Create file name with timestamp
	fileName := time.Now().Format("20060102150405") + "-" + transactionID + filepath.Ext(handler.Filename)
	filePath := filepath.Join(absolutePath, fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		return dto.AttachmentsResponse{}, errors.New("failed to create file")
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return dto.AttachmentsResponse{}, errors.New("failed to save file")
	}

	// Save attachment to database
	TransactionUUID, err := uuid.Parse(transactionID)
	if err != nil {
		return dto.AttachmentsResponse{}, errors.New("invalid transaction id")
	}

	attachment, err := transaction_serv.attachmentRepo.CreateAttachment(ctx, nil, entity.Attachments{
		Image:         fileName,
		TransactionID: TransactionUUID,
	})
	if err != nil {
		return dto.AttachmentsResponse{}, errors.New("failed to create attachment")
	}

	attachmentResponse := dto.AttachmentsResponse{
		ID:            attachment.ID.String(),
		Image:         attachment.Image,
		TransactionID: attachment.TransactionID.String(),
		CreatedAt:     attachment.CreatedAt.String(),
	}

	return attachmentResponse, nil
}

func (transaction_serv *transactionsService) UpdateTransaction(ctx context.Context, id string, transaction dto.TransactionsRequest) (dto.TransactionsResponse, error) {
	// Check if transaction exist
	transactionExist, err := transaction_serv.transactionRepo.GetTransactionByID(ctx, nil, id)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("transaction not found")
	}

	// Update transaction only if the field is not empty
	if transaction.CategoryID != "" {
		CategoryID, err := helper.ParseUUID(transaction.CategoryID)
		if err != nil {
			return dto.TransactionsResponse{}, errors.New("invalid category id")
		}

		transactionExist.CategoryID = CategoryID
	}
	if transaction.WalletID != "" {
		WalletID, err := helper.ParseUUID(transaction.WalletID)
		if err != nil {
			return dto.TransactionsResponse{}, errors.New("invalid wallet id")
		}

		transactionExist.WalletID = WalletID
	}
	if transaction.Amount != 0 {
		transactionExist.Amount = transaction.Amount
	}
	if !transaction.Date.IsZero() {
		transactionExist.TransactionDate = transaction.Date
	}
	if transaction.Description != "" {
		transactionExist.Description = transaction.Description
	}

	transactionUpdated, err := transaction_serv.transactionRepo.UpdateTransaction(ctx, nil, transactionExist)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("failed to update transaction")
	}

	transactionResponse := helper.ConvertToResponseType(transactionUpdated).(dto.TransactionsResponse)

	return transactionResponse, nil
}

func (transaction_serv *transactionsService) DeleteTransaction(ctx context.Context, id string) (dto.TransactionsResponse, error) {
	tx, err := transaction_serv.txManager.Begin(ctx)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("failed to create transaction")
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	// Check if transaction exist
	transactionExist, err := transaction_serv.transactionRepo.GetTransactionByID(ctx, tx, id)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("transaction not found")
	}

	// Update wallet balance
	if transactionExist.Category.Type == "expense" {
		transactionExist.Wallet.Balance += transactionExist.Amount
	} else if transactionExist.Category.Type == "income" {
		transactionExist.Wallet.Balance -= transactionExist.Amount
	} else {
		if transactionExist.Category.Name == "Cash Out" {
			transactionExist.Wallet.Balance += transactionExist.Amount
		} else if transactionExist.Category.Name == "Cash In" {
			transactionExist.Wallet.Balance -= transactionExist.Amount
		} else {
			return dto.TransactionsResponse{}, errors.New("invalid transaction type")
		}
	}

	// Update wallet balance
	_, err = transaction_serv.walletRepo.UpdateWallet(ctx, tx, transactionExist.Wallet)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("failed to update wallet")
	}

	// Delete transaction
	transactionDeleted, err := transaction_serv.transactionRepo.DeleteTransaction(ctx, tx, transactionExist)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("failed to delete transaction")
	}

	// Commit transaksi jika semua sukses
	if err := tx.Commit(); err != nil {
		return dto.TransactionsResponse{}, errors.New("failed to commit transaction")
	}

	transactionResponse := helper.ConvertToResponseType(transactionDeleted).(dto.TransactionsResponse)

	return transactionResponse, nil
}
