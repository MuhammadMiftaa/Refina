package handler

import (
	"net/http"

	"server/internal/entity"
	"server/internal/helper"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

type investmentHandler struct {
	investmentService service.InvestmentsService
}

func NewInvestmentHandler(investmentService service.InvestmentsService) *investmentHandler {
	return &investmentHandler{investmentService}
}

func (investment_handler *investmentHandler) GetAllInvestments(c *gin.Context) {
	investments, err := investment_handler.investmentService.GetAllInvestments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	var investmentsResponse []entity.InvestmentsResponse
	for _, investment := range investments {
		investmentResponse := helper.ConvertToResponseType(investment).(entity.InvestmentsResponse)
		investmentsResponse = append(investmentsResponse, investmentResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get All Investments",
		"data":       investmentsResponse,
	})
}

func (investment_handler *investmentHandler) GetInvestmentByID(c *gin.Context) {
	id := c.Param("id")

	investment, err := investment_handler.investmentService.GetInvestmentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	investmentResponse := helper.ConvertToResponseType(investment).(entity.InvestmentsResponse)
	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get Investment by ID",
		"data":       investmentResponse,
	})
}

func (investment_handler *investmentHandler) GetInvestmentsByUserID(c *gin.Context) {
	id := c.Param("id")

	investments, err := investment_handler.investmentService.GetInvestmentsByUserID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	var investmentsResponse []entity.InvestmentsResponse
	for _, investment := range investments {
		investmentResponse := helper.ConvertToResponseType(investment).(entity.InvestmentsResponse)
		investmentsResponse = append(investmentsResponse, investmentResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get Investments by User ID",
		"data":       investmentsResponse,
	})
}

func (investment_handler *investmentHandler) CreateInvestment(c *gin.Context) {
	var investmentRequest entity.InvestmentsRequest
	if err := c.ShouldBindJSON(&investmentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	investment, err := investment_handler.investmentService.CreateInvestment(investmentRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	investmentResponse := helper.ConvertToResponseType(investment).(entity.InvestmentsResponse)

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Create Investment",
		"data":       investmentResponse,
	})
}

func (investment_handler *investmentHandler) UpdateInvestment(c *gin.Context) {
	id := c.Param("id")

	var investmentRequest entity.InvestmentsRequest
	if err := c.ShouldBindJSON(&investmentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	investment, err := investment_handler.investmentService.UpdateInvestment(id, investmentRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	investmentResponse := helper.ConvertToResponseType(investment).(entity.InvestmentsResponse)

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Update Investment",
		"data":       investmentResponse,
	})
}

func (investment_handler *investmentHandler) DeleteInvestment(c *gin.Context) {
	id := c.Param("id")

	investment, err := investment_handler.investmentService.DeleteInvestment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	investmentResponse := helper.ConvertToResponseType(investment).(entity.InvestmentsResponse)

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Delete Investment",
		"data":       investmentResponse,
	})
}
