package middleware

import (
	"database/sql"

	"new-auth-service/internal/db"
	"new-auth-service/internal/errutil"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const PoolClientContextKey = "authenticatedAppClient"

// PoolClientCredentials validates X-Client-ID and X-Client-Secret headers
// against the app_clients table and attaches the authenticated AppClient to context.
func PoolClientCredentials() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.GetHeader("X-Client-ID")
		clientSecret := c.GetHeader("X-Client-Secret")
		if clientID == "" || clientSecret == "" {
			panic(errutil.Unauthorized("missing or malformed client credentials"))
		}

		appClient, secretHash, err := db.GetAppClientByClientID(c.Request.Context(), clientID)
		if err != nil {
			if err == sql.ErrNoRows {
				panic(errutil.Unauthorized("invalid client credentials"))
			}
			panic(errutil.Internal("failed to verify client credentials"))
		}

		if !appClient.IsActive {
			panic(errutil.Unauthorized("app client is inactive"))
		}

		if err := bcrypt.CompareHashAndPassword([]byte(secretHash), []byte(clientSecret)); err != nil {
			panic(errutil.Unauthorized("invalid client credentials"))
		}

		c.Set(PoolClientContextKey, appClient)
		c.Next()
	}
}
