package handler

import (
	"net/http"

	"server/internal/entity"
	"server/internal/helper"
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
	wallets, err := wallet_handler.walletService.GetAllWallets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	var walletsResponse []entity.WalletsResponse
	for _, wallet := range wallets {
		walletResponse := helper.ConvertToResponseType(wallet).(entity.WalletsResponse)
		walletsResponse = append(walletsResponse, walletResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get all wallets",
		"data":       walletsResponse,
	})
}

func (wallet_handler *walletHandler) GetWalletByID(c *gin.Context) {
	id := c.Param("id")

	wallet, err := wallet_handler.walletService.GetWalletByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	walletResponse := helper.ConvertToResponseType(wallet).(entity.WalletsResponse)
	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get wallet by ID",
		"data":       walletResponse,
	})
}

func (wallet_handler *walletHandler) GetWalletsByUserID(c *gin.Context) {
	id := c.Param("id")

	wallets, err := wallet_handler.walletService.GetWalletsByUserID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	var walletsResponse []entity.WalletsResponse
	for _, wallet := range wallets {
		walletResponse := helper.ConvertToResponseType(wallet).(entity.WalletsResponse)
		walletsResponse = append(walletsResponse, walletResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get wallets by user ID",
		"data":       walletsResponse,
	})
}

func (wallet_handler *walletHandler) CreateWallet(c *gin.Context) {
	var walletRequest entity.WalletsRequest
	if err := c.ShouldBindJSON(&walletRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	wallet, err := wallet_handler.walletService.CreateWallet(walletRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	walletResponse := helper.ConvertToResponseType(wallet).(entity.WalletsResponse)
	c.JSON(http.StatusCreated, gin.H{
		"statusCode": 201,
		"status":     true,
		"message":    "Create wallet",
		"data":       walletResponse,
	})
}

func (wallet_handler *walletHandler) UpdateWallet(c *gin.Context) {
	id := c.Param("id")

	var walletRequest entity.WalletsRequest
	if err := c.ShouldBindJSON(&walletRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	wallet, err := wallet_handler.walletService.UpdateWallet(id, walletRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	walletResponse := helper.ConvertToResponseType(wallet).(entity.WalletsResponse)
	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Update wallet",
		"data":       walletResponse,
	})
}

func (wallet_handler *walletHandler) DeleteWallet(c *gin.Context) {
	id := c.Param("id")

	wallet, err := wallet_handler.walletService.DeleteWallet(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	walletResponse := helper.ConvertToResponseType(wallet).(entity.WalletsResponse)
	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Delete wallet",
		"data":       walletResponse,
	})
}
