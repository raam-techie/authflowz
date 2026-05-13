package middleware

import (
	"net/http"
	"new-auth-service/utils"

	"github.com/gin-gonic/gin"
)

// TenantBearerAuth is a middleware that validates the Bearer token for tenant-level routes.
// It expects the token to be in the Authorization header in the format "Bearer <token>".
// If the token is valid, it extracts the tenant ID from the token claims and attaches it to the context.
//
// returns:
// - gin.HandlerFunc: the middleware function to be used in route groups that require tenant authentication
func TenantBearerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or malformed Authorization header"})
			return
		}

		token := authHeader[7:]

		tenant, err := utils.ParseTenantToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("tenantID", tenant.TenantID)
		c.Next()
	}
}
