package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSConfig defines the configuration for CORS middleware.
// It includes allowed origins, methods, headers, exposed headers, and max age for preflight requests.
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	ExposedHeaders []string
	MaxAge         time.Duration
}

// CORSMiddleware returns a Gin middleware function that handles CORS requests based on the provided configuration.
// The middleware sets the appropriate CORS headers and handles preflight OPTIONS requests.
//
// params:
// - cfg: CORSConfig struct containing the allowed origins, methods, headers, exposed headers, and max age for CORS requests.
//
// returns:
// - gin.HandlerFunc: the middleware function to be used in the Gin router
func CORSMiddleware(cfg CORSConfig) gin.HandlerFunc {
	origins := strings.Join(cfg.AllowedOrigins, ", ")
	methods := strings.Join(cfg.AllowedMethods, ", ")
	headers := strings.Join(cfg.AllowedHeaders, ", ")
	expose := strings.Join(cfg.ExposedHeaders, ", ")

	return func(c *gin.Context) {
		h := c.Writer.Header()

		h.Set("Access-Control-Allow-Origin", origins)
		h.Set("Access-Control-Allow-Methods", methods)
		h.Set("Access-Control-Allow-Headers", headers)

		if expose != "" {
			h.Set("Access-Control-Expose-Headers", expose)
		}

		if cfg.MaxAge > 0 {
			h.Set("Access-Control-Max-Age", strconv.Itoa(int(cfg.MaxAge.Seconds())))
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
