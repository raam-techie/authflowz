package server

import (
	"new-auth-service/internal/handlers"
	"new-auth-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerRoutes(v1 *gin.RouterGroup) {

	v1.POST("/auth/refresh", handlers.RefreshToken)

	tenant := v1.Group("/tenant")
	{
		tenant.POST("/create", handlers.CreateTenant)
		tenant.POST("/login", handlers.TenantLogin)
	}

	pool := v1.Group("/user-pool")
	{
		pool.POST("/create", handlers.CreateUserPool)
		pool.GET("/fetch", handlers.GetUserPool)
	}

	appClient := v1.Group("/app-client")
	{
		appClient.POST("/create", handlers.CreateAppClient)
		appClient.GET("/fetch", handlers.GetAppClients)
	}

	v1.GET("/user/fetch", handlers.GetUser)
	v1.GET("/audit-logs", handlers.GetAuditLogs)

	user := v1.Group("/user", middleware.PoolClientCredentials())
	{
		user.POST("/create", handlers.CreateUser)
		user.POST("/emailpassword", handlers.UserLogin)
		user.POST("/otp/email/send", handlers.SendEmailOTP)
		user.POST("/otp/email/verify", handlers.VerifyEmailOTP)
		user.PUT("/:userId/attributes", handlers.UpdateUserAttributes)
	}
}
