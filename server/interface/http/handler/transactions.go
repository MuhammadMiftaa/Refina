package handler

import (
	"net/http"

	"server/internal/dto"
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

func (transactionHandler *TransactionHandler) GetTransactionsByUserID(c *gin.Context) {
	token := c.GetHeader("Authorization")

	ctx := c.Request.Context()

	transactions, err := transactionHandler.transactionServ.GetTransactionsByUserID(ctx, token)
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

	types := c.Param("type")

	if types != "fund_transfer" {
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
	} else {
		var transaction dto.FundTransferRequest
		if err := c.ShouldBindJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":     false,
				"statusCode": 400,
				"message":    err.Error(),
			})
			return
		}
		transactionCreated, err := transactionHandler.transactionServ.FundTransfer(ctx, transaction)
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
}

func (transactionHandler *TransactionHandler) UploadAttachment(c *gin.Context) {
	ID := c.Param("id")
	ctx := c.Request.Context()

	// Check if attachment is exist
	file, handler, err := c.Request.FormFile("attachment")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}
	defer file.Close()

	attachment, err := transactionHandler.transactionServ.UploadAttachment(ctx, ID, file, handler)
	if err != nil {
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
		"data":       attachment,
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
