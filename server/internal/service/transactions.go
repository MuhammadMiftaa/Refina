package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"server/helper"
	"server/internal/repository"
	"server/internal/types/dto"
	"server/internal/types/entity"
	"server/internal/types/view"

	"github.com/google/uuid"
)

type TransactionsService interface {
	GetAllTransactions(ctx context.Context) ([]dto.TransactionsResponse, error)
	GetTransactionByID(ctx context.Context, id string) (view.ViewUserTransactions, error)
	GetTransactionsByUserID(ctx context.Context, token string) ([]view.ViewUserTransactions, error)
	CreateTransaction(ctx context.Context, transaction dto.TransactionsRequest) (dto.TransactionsResponse, error)
	FundTransfer(ctx context.Context, transaction dto.FundTransferRequest) (dto.FundTransferResponse, error)
	UploadAttachment(ctx context.Context, transactionID string, files []string) ([]dto.AttachmentsResponse, error)
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

func (transaction_serv *transactionsService) GetTransactionByID(ctx context.Context, id string) (view.ViewUserTransactions, error) {
	transaction, err := transaction_serv.transactionRepo.GetTransactionByIDJoin(ctx, nil, id)
	if err != nil {
		return view.ViewUserTransactions{}, errors.New("transaction not found")
	}

	attachments, err := transaction_serv.attachmentRepo.GetAttachmentsByTransactionID(ctx, nil, transaction.ID)
	if err != nil {
		return view.ViewUserTransactions{}, errors.New("failed to get attachments")
	}

	if len(attachments) > 0 {
		for _, attachment := range attachments {
			if attachment.Image != "" {
				result := view.Attachment{
					ID:            attachment.ID.String(),
					TransactionID: attachment.TransactionID.String(),
					Image:         attachment.Image,
					Format:        attachment.Format,
					Size:          attachment.Size,
				}
				transaction.Attachments = append(transaction.Attachments, result)
			}
		}
	} else {
		transaction.Attachments = make([]view.Attachment, 0, len(attachments))
	}

	return transaction, nil
}

func (transaction_serv *transactionsService) GetTransactionsByUserID(ctx context.Context, token string) ([]view.ViewUserTransactions, error) {
	userData, err := helper.VerifyToken(token[7:])
	if err != nil {
		return nil, errors.New("invalid token")
	}

	transactions, err := transaction_serv.transactionRepo.GetTransactionsByUserID(ctx, nil, userData.ID)
	if err != nil {
		return nil, errors.New("failed to get transactions")
	}

	return transactions, nil
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

func (transaction_serv *transactionsService) UploadAttachment(ctx context.Context, transactionID string, files []string) ([]dto.AttachmentsResponse, error) {
	var attachmentResponses []dto.AttachmentsResponse

	if transactionID == "" {
		return nil, errors.New("transaction ID is required")
	}
	if len(files) == 0 {
		return nil, errors.New("no files to upload")
	}

	for idx, file := range files {
		if file == "" {
			return nil, errors.New("file is empty")
		}

		// Ambil metadata prefix (contoh: "data:image/png;base64")
		prefixSplit := strings.SplitN(file, ",", 2)
		if len(prefixSplit) != 2 {
			return nil, errors.New("invalid base64 format")
		}
		base64Data := prefixSplit[1]

		// Decode the base64 string
		decodedFile, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return nil, errors.New("failed to decode base64 file")
		}

		if len(decodedFile) > helper.ATTACHMENT_MAX_SIZE {
			return nil, errors.New("file size exceeds 10MB")
		}

		if !helper.ATTACHMENT_EXT_ALLOWED[fmt.Sprintf(".%s", strings.Split(http.DetectContentType(decodedFile), "/")[1])] {
			return nil, errors.New("invalid file type")
		}

		// Generate file name
		mimeType := http.DetectContentType(decodedFile)
		ext := strings.ToLower("." + strings.Split(mimeType, "/")[1])
		fileName := fmt.Sprintf("%s%s", helper.GenerateFileName("TA", transactionID, strconv.Itoa(idx)), ext)

		absolutePath, _ := helper.ExpandPathAndCreateDir(helper.ATTACHMENT_FILEPATH)
		if err := helper.StorageIsExist(absolutePath); err != nil {
			return nil, errors.New("storage not found")
		}
		filePath := filepath.Join(absolutePath, fileName)

		// Simpan file ke disk
		if err := os.WriteFile(filePath, decodedFile, 0o644); err != nil {
			return nil, errors.New("failed to save file")
		}

		// Save attachment to database
		TransactionUUID, err := uuid.Parse(transactionID)
		if err != nil {
			return nil, errors.New("invalid transaction id")
		}

		attachment, err := transaction_serv.attachmentRepo.CreateAttachment(ctx, nil, entity.Attachments{
			Image:         fileName,
			TransactionID: TransactionUUID,
			Size:          int64(len(decodedFile)),
			Format:        ext,
		})
		if err != nil {
			return nil, errors.New("failed to create attachment")
		}

		attachmentResponse := dto.AttachmentsResponse{
			ID:            attachment.ID.String(),
			Image:         attachment.Image,
			TransactionID: attachment.TransactionID.String(),
			CreatedAt:     attachment.CreatedAt.String(),
		}

		attachmentResponses = append(attachmentResponses, attachmentResponse)
	}

	return attachmentResponses, nil
}

func (transaction_serv *transactionsService) UpdateTransaction(ctx context.Context, id string, transaction dto.TransactionsRequest) (dto.TransactionsResponse, error) {
	// ! Begin a new transaction
	tx, err := transaction_serv.txManager.Begin(ctx)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("failed to create transaction")
	}

	// ! Defer rollback if there is an error or panic
	defer func() {
		// Rollback otomatis jika transaksi belum di-commit
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	// ? Check if transaction exist
	transactionExist, err := transaction_serv.transactionRepo.GetTransactionByID(ctx, tx, id)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("transaction not found")
	}

	// ? If category ID is different, update category
	if transaction.CategoryID != transactionExist.CategoryID.String() {
		CategoryID, err := helper.ParseUUID(transaction.CategoryID)
		if err != nil {
			return dto.TransactionsResponse{}, errors.New("invalid category id")
		}

		// * Check if category exist
		_, err = transaction_serv.categoryRepo.GetCategoryByID(ctx, tx, transaction.CategoryID)
		if err != nil {
			return dto.TransactionsResponse{}, errors.New("category not found")
		}

		transactionExist.CategoryID = CategoryID
	}

	// ? If wallet ID is different, update wallet balance
	if transaction.WalletID != transactionExist.WalletID.String() {
		// *  Check if wallet exist
		oldWallet, err := transaction_serv.walletRepo.GetWalletByID(ctx, tx, transactionExist.WalletID.String())
		if err != nil {
			return dto.TransactionsResponse{}, errors.New("wallet not found")
		}

		// *  Update wallet balance
		switch transactionExist.Category.Type {
		case "expense":
			oldWallet.Balance += transactionExist.Amount
		case "income":
			oldWallet.Balance -= transactionExist.Amount
		default:
			return dto.TransactionsResponse{}, errors.New("invalid transaction type")
		}

		if _, err = transaction_serv.walletRepo.UpdateWallet(ctx, tx, oldWallet); err != nil {
			return dto.TransactionsResponse{}, errors.New("failed to update wallet")
		}

		// *  Check if new wallet exist
		newWallet, err := transaction_serv.walletRepo.GetWalletByID(ctx, tx, transaction.WalletID)
		if err != nil {
			return dto.TransactionsResponse{}, errors.New("new wallet not found")
		}

		// *  Update wallet balance
		switch transactionExist.Category.Type {
		case "expense":
			newWallet.Balance -= transaction.Amount
		case "income":
			newWallet.Balance += transaction.Amount
		default:
			return dto.TransactionsResponse{}, errors.New("invalid transaction type")
		}

		if _, err = transaction_serv.walletRepo.UpdateWallet(ctx, tx, newWallet); err != nil {
			return dto.TransactionsResponse{}, errors.New("failed to update new wallet")
		}

		// *  Parse ID from JSON to valid UUID
		WalletID, err := helper.ParseUUID(transaction.WalletID)
		if err != nil {
			return dto.TransactionsResponse{}, errors.New("invalid wallet id")
		}
		transactionExist.WalletID = WalletID
	}

	// ? Update transaction fields
	if transaction.Amount != transactionExist.Amount {
		// *  Update wallet balance
		oldWallet, err := transaction_serv.walletRepo.GetWalletByID(ctx, tx, transactionExist.WalletID.String())
		if err != nil {
			return dto.TransactionsResponse{}, errors.New("wallet not found")
		}

		// *  Update wallet balance
		switch transactionExist.Category.Type {
		case "expense":
			oldWallet.Balance += transactionExist.Amount
			oldWallet.Balance -= transaction.Amount
		case "income":
			oldWallet.Balance -= transactionExist.Amount
			oldWallet.Balance += transaction.Amount
		default:
			return dto.TransactionsResponse{}, errors.New("invalid transaction type")
		}

		if _, err = transaction_serv.walletRepo.UpdateWallet(ctx, tx, oldWallet); err != nil {
			return dto.TransactionsResponse{}, errors.New("failed to update wallet")
		}

		// *  Update transaction amount
		transactionExist.Amount = transaction.Amount
	}

	// ? Update transaction date
	if !transaction.Date.IsZero() {
		transactionExist.TransactionDate = transaction.Date
	}

	// ? Update description
	if transaction.Description != "" {
		transactionExist.Description = transaction.Description
	}

	// ? Update transaction
	transactionUpdated, err := transaction_serv.transactionRepo.UpdateTransaction(ctx, tx, transactionExist)
	if err != nil {
		return dto.TransactionsResponse{}, errors.New("failed to update transaction")
	}

	// ? If attachments exist, update attachments
	if len(transaction.Attachments) > 0 {
		for _, attachment := range transaction.Attachments {
			switch attachment.Status {
			case "create":
				// * Create new attachment
				if len(attachment.Files) == 0 {
					return dto.TransactionsResponse{}, errors.New("no files to upload")
				}

				if _, err := transaction_serv.UploadAttachment(ctx, transactionUpdated.ID.String(), attachment.Files); err != nil {
					return dto.TransactionsResponse{}, fmt.Errorf("failed to upload attachment: %w", err)
				}

			case "delete":
				// * Delete attachment
				if len(attachment.Files) == 0 {
					return dto.TransactionsResponse{}, errors.New("no files to delete")
				}

				for _, ID := range attachment.Files {
					// * Get attachment by ID
					attachmentToDelete, err := transaction_serv.attachmentRepo.GetAttachmentByID(ctx, tx, ID)
					if err != nil {
						return dto.TransactionsResponse{}, fmt.Errorf("attachment with file %s not found: %w", ID, err)
					}

					// * Check if attachment belongs to transaction
					if attachmentToDelete.TransactionID != transactionUpdated.ID {
						return dto.TransactionsResponse{}, fmt.Errorf("attachment with file %s does not belong to transaction %s", ID, transactionUpdated.ID)
					}

					// * Delete file from database
					if _, err := transaction_serv.attachmentRepo.DeleteAttachment(ctx, tx, attachmentToDelete); err != nil {
						return dto.TransactionsResponse{}, fmt.Errorf("attachment with file %v not found: %w", attachmentToDelete, err)
					}
				}
			
			default:
				return dto.TransactionsResponse{}, errors.New("invalid attachment status")
			}
		}
	}

	// ! Commit transaction if all operations are successful
	if err = tx.Commit(); err != nil {
		return dto.TransactionsResponse{}, errors.New("failed to commit transaction")
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
