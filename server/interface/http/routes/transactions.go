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
	txManager := repository.NewTxManager(db)
	transactionRepo := repository.NewTransactionRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	attachmentRepo := repository.NewAttachmentsRepository(db)

	Transaction_serv := service.NewTransactionService(txManager, transactionRepo, walletRepo, categoryRepo, attachmentRepo)
	Transaction_handler := handler.NewTransactionHandler(Transaction_serv)

	version.Use(middleware.AuthMiddleware())
	version.GET("transactions", Transaction_handler.GetAllTransactions)
	version.GET("transactions/:id", Transaction_handler.GetTransactionByID)
	version.GET("transactions/user", Transaction_handler.GetTransactionsByUserID)
	version.POST("transactions/:type", Transaction_handler.CreateTransaction)
	version.POST("transactions/attachment/:id", Transaction_handler.UploadAttachment)
	version.PUT("transactions/:id", Transaction_handler.UpdateTransaction)
	version.DELETE("transactions/:id", Transaction_handler.DeleteTransaction)
}
