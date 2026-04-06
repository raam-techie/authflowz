package server

import (
	"new-auth-service/routes"

	"github.com/gin-gonic/gin"
)

func registerRoutes(v1 *gin.RouterGroup) {
	routes.RegisterTenantRoutes(v1)
	routes.RegisterClientRoutes(v1)
	routes.RegisterUserRoutes(v1)
}
