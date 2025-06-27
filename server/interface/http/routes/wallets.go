package routes

import (
	"server/interface/http/handler"
	"server/interface/http/middleware"
	"server/internal/repository"
	"server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WalletRoutes(version *gin.RouterGroup, db *gorm.DB) {
	txManager := repository.NewTxManager(db)
	Wallet_repo := repository.NewWalletRepository(db)
	Wallet_serv := service.NewWalletService(txManager, Wallet_repo)
	Wallet_handler := handler.NewWalletHandler(Wallet_serv)

	wallets := version.Group("/wallets")
	wallets.Use(middleware.AuthMiddleware())

	wallets.GET("", Wallet_handler.GetAllWallets)
	wallets.GET(":id", Wallet_handler.GetWalletByID)
	wallets.GET("user", Wallet_handler.GetWalletsByUserID)
	wallets.POST("", Wallet_handler.CreateWallet)
	wallets.PUT(":id", Wallet_handler.UpdateWallet)
	wallets.DELETE(":id", Wallet_handler.DeleteWallet)
}
