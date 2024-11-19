package handler

import (
	"errors"
	"net/http"
	"path/filepath"
	"time"

	"server/internal/entity"
	"server/internal/helper"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionServ service.TransactionsService
}

func NewTransactionHandler(transactionServ service.TransactionsService) *TransactionHandler {
	return &TransactionHandler{transactionServ}
}

func (transactionHandler *TransactionHandler) GetAllTransactions(c *gin.Context) {
	transactions, err := transactionHandler.transactionServ.GetAllTransactions()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	var transactionsResponse []entity.TransactionsResponse
	for _, transaction := range transactions {
		transactionResponse := helper.ConvertToResponseType(transaction).(entity.TransactionsResponse)
		transactionsResponse = append(transactionsResponse, transactionResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "Get all transactions data",
		"data":       transactionsResponse,
	})
}

func (transactionHandler *TransactionHandler) GetTransactionByID(c *gin.Context) {
	id := c.Param("id")

	transaction, err := transactionHandler.transactionServ.GetTransactionByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	transactionResponse := helper.ConvertToResponseType(transaction).(entity.TransactionsResponse)

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "Get transaction data by ID",
		"data":       transactionResponse,
	})
}

func (transactionHandler *TransactionHandler) CreateTransaction(c *gin.Context) {
	var transaction entity.TransactionsRequest
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	transactionCreated, err := transactionHandler.transactionServ.CreateTransaction(transaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	transactionResponse := helper.ConvertToResponseType(transactionCreated).(entity.TransactionsResponse)

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "Create transaction data",
		"data":       transactionResponse,
	})
}

func (transactionHandler *TransactionHandler) UploadAttachment(c *gin.Context) {
	path := "../client/src/assets/attachments"
	absoultePath, _ := filepath.Abs(path)

	// Check if attachment is exist
	file, err := c.FormFile("attachment")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	// Check if storage is exist
	if err = helper.StorageIsExist(absoultePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     false,
			"statusCode": 500,
			"message":    errors.New("Storage not found").Error(),
		})
		return
	}

	// Create file name with timestamp
	fileName := time.Now().Format("20060102150405") + filepath.Ext(file.Filename)
	filePath := filepath.Join(absoultePath, fileName)
	if err = c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     false,
			"statusCode": 500,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "Upload attachment success",
		"data":       filePath,
	})
}

func (transactionHandler *TransactionHandler) UpdateTransaction(c *gin.Context) {
	var transaction entity.TransactionsRequest
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	id := c.Param("id")

	transactionUpdated, err := transactionHandler.transactionServ.UpdateTransaction(id, transaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	transactionResponse := helper.ConvertToResponseType(transactionUpdated).(entity.TransactionsResponse)

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "Update transaction data",
		"data":       transactionResponse,
	})
}

func (transactionHandler *TransactionHandler) DeleteTransaction(c *gin.Context) {
	id := c.Param("id")

	transactionDeleted, err := transactionHandler.transactionServ.DeleteTransaction(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	transactionResponse := helper.ConvertToResponseType(transactionDeleted).(entity.TransactionsResponse)

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "Delete transaction data",
		"data":       transactionResponse,
	})
}
