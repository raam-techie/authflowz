package routes

import (
	"github.com/gin-gonic/gin"
	"new-auth-service/handlers"
)

func RegisterTenantRoutes(v1 *gin.RouterGroup) {

	
	tenant := v1.Group("/tenant")
	{
		tenant.POST("/create", handlers.CreateTenant)
	}
}
