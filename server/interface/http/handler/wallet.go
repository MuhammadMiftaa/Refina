package handler

import (
	"net/http"

	"server/internal/dto"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

type walletHandler struct {
	walletService service.WalletsService
}

func NewWalletHandler(walletService service.WalletsService) *walletHandler {
	return &walletHandler{walletService}
}

func (wallet_handler *walletHandler) GetAllWallets(c *gin.Context) {
	ctx := c.Request.Context()

	wallets, err := wallet_handler.walletService.GetAllWallets(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get all wallets",
		"data":       wallets,
	})
}

func (wallet_handler *walletHandler) GetWalletByID(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	wallet, err := wallet_handler.walletService.GetWalletByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get wallet by ID",
		"data":       wallet,
	})
}

func (wallet_handler *walletHandler) GetWalletsByUserID(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	wallets, err := wallet_handler.walletService.GetWalletsByUserID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get wallets by user ID",
		"data":       wallets,
	})
}

func (wallet_handler *walletHandler) CreateWallet(c *gin.Context) {
	ctx := c.Request.Context()

	var walletRequest dto.WalletsRequest
	if err := c.ShouldBindJSON(&walletRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	wallet, err := wallet_handler.walletService.CreateWallet(ctx, walletRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": 201,
		"status":     true,
		"message":    "Create wallet",
		"data":       wallet,
	})
}

func (wallet_handler *walletHandler) UpdateWallet(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	var walletRequest dto.WalletsRequest
	if err := c.ShouldBindJSON(&walletRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	wallet, err := wallet_handler.walletService.UpdateWallet(ctx, id, walletRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Update wallet",
		"data":       wallet,
	})
}

func (wallet_handler *walletHandler) DeleteWallet(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	wallet, err := wallet_handler.walletService.DeleteWallet(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Delete wallet",
		"data":       wallet,
	})
}
