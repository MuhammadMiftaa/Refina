package router

import (
	"server/interface/http/middleware"
	"server/interface/http/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	v1 := router.Group("/v1")
	routes.UserRoutes(v1, db)
	routes.TransactionRoutes(v1, db)

	return router
}