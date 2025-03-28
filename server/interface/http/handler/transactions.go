package handler

import (
	"errors"
	"net/http"
	"path/filepath"
	"time"

	"server/internal/dto"
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
	ctx := c.Request.Context()

	transactions, err := transactionHandler.transactionServ.GetAllTransactions(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "Get all transactions data",
		"data":       transactions,
	})
}

func (transactionHandler *TransactionHandler) GetTransactionByID(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	transaction, err := transactionHandler.transactionServ.GetTransactionByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "Get transaction data by ID",
		"data":       transaction,
	})
}

func (transactionHandler *TransactionHandler) GetTransactionsByWalletID(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	transactions, err := transactionHandler.transactionServ.GetTransactionsByWalletID(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "Get transactions data by wallet ID",
		"data":       transactions,
	})
}

func (transactionHandler *TransactionHandler) CreateTransaction(c *gin.Context) {
	ctx := c.Request.Context()

	var transaction dto.TransactionsRequest
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	transactionCreated, err := transactionHandler.transactionServ.CreateTransaction(ctx, transaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "Create transaction data",
		"data":       transactionCreated,
	})
}

func (transactionHandler *TransactionHandler) UploadAttachment(c *gin.Context) {
	path := "../refina/src/assets/attachments"
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
			"message":    errors.New("storage not found").Error(),
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
	ctx := c.Request.Context()

	var transaction dto.TransactionsRequest
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	id := c.Param("id")

	transactionUpdated, err := transactionHandler.transactionServ.UpdateTransaction(ctx, id, transaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "Update transaction data",
		"data":       transactionUpdated,
	})
}

func (transactionHandler *TransactionHandler) DeleteTransaction(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	transactionDeleted, err := transactionHandler.transactionServ.DeleteTransaction(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "Delete transaction data",
		"data":       transactionDeleted,
	})
}
