package routes

import (
	"server/interface/http/handler"
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
	
	// version.Use(middleware.AuthMiddleware())
	version.GET("investments", Investment_handler.GetAllInvestments)
	version.GET("investments/:id", Investment_handler.GetInvestmentByID)
	version.GET("investments/user/:id", Investment_handler.GetInvestmentsByUserID)
	version.POST("investments", Investment_handler.CreateInvestment)
	version.PUT("investments/:id", Investment_handler.UpdateInvestment)
	version.DELETE("investments/:id", Investment_handler.DeleteInvestment)
}
