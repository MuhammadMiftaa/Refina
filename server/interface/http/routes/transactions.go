package routes

import (
	"server/interface/http/handler"
	"server/interface/http/middleware"
	"server/internal/repository"
	"server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TransactionRoutes(version *gin.RouterGroup, db *gorm.DB) {
	Transaction_repo := repository.NewTransactionRepository(db)
	User_repo := repository.NewUsersRepository(db)
	Transaction_serv := service.NewTransactionService(Transaction_repo, User_repo)
	Transaction_handler := handler.NewTransactionHandler(Transaction_serv)

	version.Use(middleware.AuthMiddleware())
	version.GET("transactions", Transaction_handler.GetAllTransactions)
	version.GET("transactions/:id", Transaction_handler.GetTransactionByID)
	version.POST("transactions", Transaction_handler.CreateTransaction)
	version.PUT("transactions/:id", Transaction_handler.UpdateTransaction)
	version.DELETE("transactions/:id", Transaction_handler.DeleteTransaction)
}
