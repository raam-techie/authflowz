package server

import (
	"new-auth-service/db"
	"new-auth-service/errutil"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Client-ID", "X-Client-Secret"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
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
