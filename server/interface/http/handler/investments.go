package handler

import (
	"net/http"

	"server/internal/service"
	"server/internal/types/dto"

	"github.com/gin-gonic/gin"
)

type investmentHandler struct {
	investmentService service.InvestmentsService
}

func NewInvestmentHandler(investmentService service.InvestmentsService) *investmentHandler {
	return &investmentHandler{investmentService}
}

func (investment_handler *investmentHandler) GetAllInvestments(c *gin.Context) {
	ctx := c.Request.Context()

	investments, err := investment_handler.investmentService.GetAllInvestments(ctx)
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
		"message":    "Get All Investments",
		"data":       investments,
	})
}

func (investment_handler *investmentHandler) GetInvestmentByID(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	investment, err := investment_handler.investmentService.GetInvestmentByID(ctx, id)
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
		"message":    "Get Investment by ID",
		"data":       investment,
	})
}

func (investment_handler *investmentHandler) GetInvestmentsByUserID(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.GetHeader("Authorization")

	userWallets, err := investment_handler.investmentService.GetInvestmentByID(ctx, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get user wallets data",
		"data":       userWallets,
	})
}

func (investment_handler *investmentHandler) CreateInvestment(c *gin.Context) {
	ctx := c.Request.Context()

	var investmentRequest dto.InvestmentsRequest
	if err := c.ShouldBindJSON(&investmentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	investment, err := investment_handler.investmentService.CreateInvestment(ctx, investmentRequest)
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
		"message":    "Create Investment",
		"data":       investment,
	})
}

func (investment_handler *investmentHandler) UpdateInvestment(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	var investmentRequest dto.InvestmentsRequest
	if err := c.ShouldBindJSON(&investmentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	investment, err := investment_handler.investmentService.UpdateInvestment(ctx, id, investmentRequest)
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
		"message":    "Update Investment",
		"data":       investment,
	})
}

func (investment_handler *investmentHandler) DeleteInvestment(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	investment, err := investment_handler.investmentService.DeleteInvestment(ctx, id)
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
		"message":    "Delete Investment",
		"data":       investment,
	})
}
