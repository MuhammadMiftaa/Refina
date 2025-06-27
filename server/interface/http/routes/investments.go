package routes

import (
	"server/interface/http/handler"
	"server/interface/http/middleware"
	"server/internal/repository"
	"server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InvestmentRoute(version *gin.RouterGroup, db *gorm.DB) {
	txManager := repository.NewTxManager(db)
	Investment_repo := repository.NewInvestmentRepository(db)
	Investment_serv := service.NewInvestmentService(txManager, Investment_repo)
	Investment_handler := handler.NewInvestmentHandler(Investment_serv)

	investments := version.Group("/investments")
	investments.Use(middleware.AuthMiddleware())

	investments.GET("", Investment_handler.GetAllInvestments)
	investments.GET(":id", Investment_handler.GetInvestmentByID)
	investments.GET("user", Investment_handler.GetInvestmentsByUserID)
	investments.POST("", Investment_handler.CreateInvestment)
	investments.PUT(":id", Investment_handler.UpdateInvestment)
	investments.DELETE(":id", Investment_handler.DeleteInvestment)
}
