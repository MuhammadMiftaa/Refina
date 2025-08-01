package router

import (
	"server/helper"
	"server/interface/http/middleware"
	"server/interface/http/routes"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, redis *redis.Client) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	router.GET("test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	v1 := router.Group("/v1")
	routes.UserRoutes(v1, db, redis)
	routes.TransactionRoutes(v1, db)
	routes.WalletRoutes(v1, db)
	routes.InvestmentRoute(v1, db)
	routes.WalletTypesRoutes(v1, db)
	routes.CategoryRoutes(v1, db)

	router.Static("/uploads/transactions-attachments", helper.ATTACHMENT_FILEPATH)

	return router
}
