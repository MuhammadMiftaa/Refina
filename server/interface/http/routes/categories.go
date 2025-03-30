package routes

import (
	"server/interface/http/handler"
	"server/internal/repository"
	"server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CategoryRoutes(version *gin.RouterGroup, db *gorm.DB) {
	txManager := repository.NewTxManager(db)
	categoryRepo := repository.NewCategoryRepository(db)

	categoryServ := service.NewCategoriesService(txManager, categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryServ)

	version.GET("categories", categoryHandler.GetAllCategories)
	version.GET("categories/:id", categoryHandler.GetCategoryByID)
	version.GET("categories/type/:type", categoryHandler.GetCategoriesByType)
	version.POST("categories", categoryHandler.CreateCategory)
	version.PUT("categories/:id", categoryHandler.UpdateCategory)
	version.DELETE("categories/:id", categoryHandler.DeleteCategory)
}
