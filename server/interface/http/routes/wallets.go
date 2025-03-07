package routes

import (
	"server/interface/http/handler"
	// "server/interface/http/middleware"
	"server/internal/repository"
	"server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WalletRoutes(version *gin.RouterGroup, db *gorm.DB) {
	Wallet_repo := repository.NewWalletRepository(db)
	Wallet_serv := service.NewWalletService(Wallet_repo)
	Wallet_handler := handler.NewWalletHandler(Wallet_serv)

	// version.Use(middleware.AuthMiddleware())
	version.GET("wallets", Wallet_handler.GetAllWallets)
	version.GET("wallets/:id", Wallet_handler.GetWalletByID)
	version.GET("wallets/user/:id", Wallet_handler.GetWalletsByUserID)
	version.POST("wallets", Wallet_handler.CreateWallet)
	version.PUT("wallets/:id", Wallet_handler.UpdateWallet)
	version.DELETE("wallets/:id", Wallet_handler.DeleteWallet)
}
