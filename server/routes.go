package server

import (
	"new-auth-service/handlers"
	"new-auth-service/middleware"

	"github.com/gin-gonic/gin"
)

func registerRoutes(v1 *gin.RouterGroup) {

	v1.POST("/auth/refresh", handlers.RefreshToken)

	tenant := v1.Group("/tenant")
	{
		tenant.POST("/create", handlers.CreateTenant)
		tenant.POST("/login", handlers.TenantLogin)
	}

	project := v1.Group("/project")
	{
		project.POST("/create", handlers.CreateProject)
		project.GET("", handlers.GetProject)
	}

	user := v1.Group("/user", middleware.ClientCredentials())
	{
		user.POST("/create", handlers.CreateUser)
		user.POST("/emailpassword", handlers.UserLogin)
		user.POST("/otp/email/send", handlers.SendEmailOTP)
		user.POST("/otp/email/verify", handlers.VerifyEmailOTP)
		user.PUT("/:userId/attributes", handlers.UpdateUserAttributes)
	}
}
