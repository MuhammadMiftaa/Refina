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

	transaction := version.Group("/transactions")
	transaction.Use(middleware.AuthMiddleware())

	transaction.GET("", Transaction_handler.GetAllTransactions)
	transaction.GET(":id", Transaction_handler.GetTransactionByID)
	transaction.GET("user", Transaction_handler.GetTransactionsByUserID)
	transaction.POST(":type", Transaction_handler.CreateTransaction)
	transaction.POST("attachment/:id", Transaction_handler.UploadAttachment)
	transaction.PUT(":id", Transaction_handler.UpdateTransaction)
	transaction.DELETE(":id", Transaction_handler.DeleteTransaction)
	transaction.GET("user-summary", Transaction_handler.GetUserSummary)
	transaction.GET("user-summary/detail", Transaction_handler.GetUserSummaryByUserID)
	transaction.GET("user-monthly-summary", Transaction_handler.GetUserMonthlySummary)
	transaction.GET("user-monthly-summary/detail", Transaction_handler.GetUserMonthlySummaryByUserID)
	transaction.GET("user-most-expenses", Transaction_handler.GetUserMostExpenses)
	transaction.GET("user-most-expenses/detail", Transaction_handler.GetUserMostExpensesByUserID)
	transaction.GET("user-wallet-daily-summary", Transaction_handler.GetUserWalletDailySummary)
	transaction.GET("user-wallet-daily-summary/detail", Transaction_handler.GetUserWalletDailySummaryByUserID)
}
