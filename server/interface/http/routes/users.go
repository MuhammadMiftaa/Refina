package routes

import (
	"server/interface/http/handler"
	"server/interface/http/middleware"
	"server/internal/repository"
	"server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func UserRoutes(version *gin.RouterGroup, db *gorm.DB, redis *redis.Client) {
	User_repo := repository.NewUsersRepository(db)
	User_serv := service.NewUsersService(User_repo)

	OTP_repo := repository.NewOTPRepository(redis)
	OTP_serv := service.NewOTPService(OTP_repo)

	User_handler := handler.NewUsersHandler(User_serv, OTP_serv)

	auth := version.Group("/auth")
	{
		auth.POST("login", User_handler.Login)
		auth.POST("register", User_handler.Register)
		auth.POST("send/otp", User_handler.SendOTP)
		auth.POST("verify/otp", User_handler.VerifyOTP)

		auth.GET("google/oauth", User_handler.OAuthHandler("google"))
		auth.GET("callback/google", User_handler.CallbackGoogle)
		auth.GET("github/oauth", User_handler.OAuthHandler("github"))
		auth.GET("callback/github", User_handler.CallbackGithub)
		auth.GET("microsoft/oauth", User_handler.OAuthHandler("microsoft"))
		auth.GET("callback/microsoft", User_handler.CallbackMicrosoft)
	}

	version.Use(middleware.AuthMiddleware())
	version.GET("users", User_handler.GetAllUsers)
	version.GET("users/:id", User_handler.GetUserByID)
	version.PUT("users/:id", User_handler.UpdateUser)
	version.DELETE("users/:id", User_handler.DeleteUser)

	version.GET("users/wallets", User_handler.GetUserWallets)
	version.GET("users/:id/investments", User_handler.GetUserInvestments)
	version.GET("users/transactions", User_handler.GetUserTransactions)
}