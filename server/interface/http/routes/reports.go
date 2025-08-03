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
	reportsRepo := repository.NewReportsRepository(db)
	reportsService := service.NewReportsService(reportsRepo, userRepo)
	reportHandler := handler.NewReportHandler(reportsService)

	reportGroup := version.Group("/reports")

	reportGroup.POST("request", reportHandler.RequestReports)
}
