package server

import (
	"new-auth-service/db"
	"new-auth-service/errutil"
	"new-auth-service/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware(middleware.CORSConfig{
		AllowedOrigins: []string{"*"}, // TODO: Restrict this in production to specific domains
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Content-Type", "Authorization", "X-Client-ID", "X-Client-Secret"}, // Extend as needed
		ExposedHeaders: []string{"Content-Length"},
		MaxAge:         12 * time.Hour,
	}))

	router.Use(errutil.ErrorMiddleware())

	router.GET("/health", func(c *gin.Context) {
		if err := db.DB.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "error", "db": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok", "db": "connected"})
	})

	v1 := router.Group("/api/v1")
	registerRoutes(v1)

	return router
}
