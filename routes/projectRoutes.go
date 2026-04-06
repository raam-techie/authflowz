package routes

import (
	"new-auth-service/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterProjectRoutes(v1 *gin.RouterGroup) {
	project := v1.Group("/project")
	{
		project.POST("/create", handlers.CreateProject)
	}
}
