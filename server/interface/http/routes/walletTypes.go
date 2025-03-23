package routes

import (
	"server/interface/http/handler"
	"server/interface/http/middleware"
	"server/internal/repository"
	"server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WalletTypesRoutes(version *gin.RouterGroup, db *gorm.DB) {
	txManager := repository.NewTxManager(db)
	WalletTypes_repo := repository.NewWalletTypesRepository(db)
	WalletTypes_serv := service.NewWalletTypesService(txManager, WalletTypes_repo)
	WalletTypes_handler := handler.NewWalletTypesHandler(WalletTypes_serv)

	version.Use(middleware.AuthMiddleware())
	version.GET("wallet-types", WalletTypes_handler.GetAllWalletTypes)
	version.GET("wallet-types/:id", WalletTypes_handler.GetWalletTypeByID)
	version.POST("wallet-types", WalletTypes_handler.CreateWalletType)
	version.PUT("wallet-types/:id", WalletTypes_handler.UpdateWalletType)
	version.DELETE("wallet-types/:id", WalletTypes_handler.DeleteWalletType)
}
