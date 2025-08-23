package router

import (
	"server/config/db"
	"server/config/miniofs"
	"server/config/redis"
	"server/interface/http/middleware"
	"server/interface/http/routes"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware(), middleware.GinMiddleware())

	router.GET("test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	v1 := router.Group("/v1")
	routes.UserRoutes(v1, db.DB, redis.RDB)
	routes.TransactionRoutes(v1, db.DB, miniofs.MinioClient)
	routes.WalletRoutes(v1, db.DB)
	routes.InvestmentRoute(v1, db.DB)
	routes.WalletTypesRoutes(v1, db.DB)
	routes.CategoryRoutes(v1, db.DB)
	routes.ReportRoutes(v1, db.DB)

	return router
}
