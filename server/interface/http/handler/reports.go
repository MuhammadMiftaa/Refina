package handler

import (
	"net/http"

	"server/internal/service"
	"server/internal/types/dto"
	"server/queue/config"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
)

type reportHandler struct {
	userService service.UsersService
}

func NewReportHandler(userService service.UsersService) *reportHandler {
	return &reportHandler{
		userService: userService,
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

	if _, err = h.userService.GetUserByID(request.UserID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":     false,
			"statusCode": 404,
			"message":    err.Error(),
		})
		return
	}

	prodRequestReports, err := config.GetChannel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     false,
			"statusCode": 500,
			"message":    "Failed to get RabbitMQ channel",
		})
		return
	}

	message := amqp091.Publishing{
		ContentType: "application/json",
		Body:        []byte(`{"user_id":"` + request.UserID + `", "from_date":"` + request.FromDate + `", "to_date":"` + request.ToDate + `"}`),
	}
	ctx := c.Request.Context()
	err = prodRequestReports.PublishWithContext(ctx, "reports", "request", false, false, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":     false,
			"statusCode": 500,
			"message":    "Failed to publish request to RabbitMQ",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "Request reports sent successfully",
	})
}
