package routes

import (
	"new-auth-service/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterClientRoutes(v1 *gin.RouterGroup) {
	client := v1.Group("/client")
	{
		client.POST("/create", handlers.CreateClient)
	}
}
