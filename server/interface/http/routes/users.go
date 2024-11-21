package routes

import (
	"server/interface/http/handler"
	"server/interface/http/middleware"
	"server/internal/repository"
	"server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(version *gin.RouterGroup, db *gorm.DB) {
	User_repo := repository.NewUsersRepository(db)
	User_serv := service.NewUsersService(User_repo)
	User_handler := handler.NewUsersHandler(User_serv)

	auth := version.Group("/auth")
	{
		auth.POST("register", User_handler.Register)
		auth.POST("login", User_handler.Login)
		auth.GET("google/oauth", User_handler.OAuthGoogle)
		auth.GET("callback/google", User_handler.CallbackGoogle)
		auth.GET("github/oauth", User_handler.OAuthGithub)
		auth.GET("callback/github", User_handler.CallbackGithub)
		auth.GET("microsoft/oauth", User_handler.OAuthMicrosoft)
		auth.GET("callback/microsoft", User_handler.CallbackMicrosoft)
	}

	version.Use(middleware.AuthMiddleware())
	version.GET("users", User_handler.GetAllUsers)
	version.GET("users/:id", User_handler.GetUserByID)
	version.PUT("users/:id", User_handler.UpdateUser)
	version.DELETE("users/:id", User_handler.DeleteUser)
}