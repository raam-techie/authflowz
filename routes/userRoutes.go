package routes

import (
	"new-auth-service/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(v1 *gin.RouterGroup) {
	user := v1.Group("/user")
	{
		user.POST("/create", handlers.CreateUser)
	}
}
