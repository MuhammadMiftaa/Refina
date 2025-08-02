package routes

import (
	"server/interface/http/handler"
	"server/internal/repository"
	"server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ReportRoutes(version *gin.RouterGroup, db *gorm.DB) {
	userRepo := repository.NewUsersRepository(db)
	UserService := service.NewUsersService(userRepo)
	reportHandler := handler.NewReportHandler(UserService)

	reportGroup := version.Group("/reports")

	reportGroup.GET("request", reportHandler.RequestReports)
}
