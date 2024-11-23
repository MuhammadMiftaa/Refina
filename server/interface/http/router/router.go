package router

import (
	"server/interface/http/middleware"
	"server/interface/http/routes"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, redis *redis.Client) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	v1 := router.Group("/v1")
	routes.UserRoutes(v1, db, redis)
	routes.TransactionRoutes(v1, db)

	return router
}