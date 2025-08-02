package router

import (
	"server/helper"
	"server/interface/http/middleware"
	"server/interface/http/routes"
	dc "server/db/config"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	router.GET("test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	v1 := router.Group("/v1")
	routes.UserRoutes(v1, dc.DB, dc.RDB)
	routes.TransactionRoutes(v1, dc.DB)
	routes.WalletRoutes(v1, dc.DB)
	routes.InvestmentRoute(v1, dc.DB)
	routes.WalletTypesRoutes(v1, dc.DB)
	routes.CategoryRoutes(v1, dc.DB)
	routes.ReportRoutes(v1, dc.DB)

	router.Static("/uploads/transactions-attachments", helper.ATTACHMENT_FILEPATH)

	return router
}
