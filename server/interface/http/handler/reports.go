package handler

import (
	"net/http"

	"server/internal/service"
	"server/internal/types/dto"

	"github.com/gin-gonic/gin"
)

type reportHandler struct {
	reportService service.ReportsService
}

func NewReportHandler(reportService service.ReportsService) *reportHandler {
	return &reportHandler{
		reportService: reportService,
	}
}

func (h *reportHandler) RequestReports(c *gin.Context) {
	var request dto.ReportRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"statusCode": 400,
			"message":    err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	if err := h.reportService.RequestReport(ctx, request.UserID, request.FromDate, request.ToDate); err != nil {
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
		"message":    "Request reports sent successfully",
	})
}
