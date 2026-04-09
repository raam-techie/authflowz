package middleware

import (
	"database/sql"

	"new-auth-service/db"
	"new-auth-service/errutil"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const ProjectContextKey = "authenticatedProject"

func ClientCredentials() gin.HandlerFunc {
	return func(c *gin.Context) {

		clientID := c.GetHeader("X-Client-ID")
		clientSecret := c.GetHeader("X-Client-Secret")
		if clientID == "" || clientSecret == "" {
			panic(errutil.Unauthorized("missing or malformed client credentials"))
		}

		project, secretHash, err := db.GetProjectByClientID(c.Request.Context(), clientID)
		if err != nil {
			if err == sql.ErrNoRows {
				panic(errutil.Unauthorized("invalid client credentials"))
			}
			panic(errutil.Internal("failed to verify client credentials"))
		}

		if !project.IsActive {
			panic(errutil.Unauthorized("project is inactive"))
		}

		if err := bcrypt.CompareHashAndPassword([]byte(secretHash), []byte(clientSecret)); err != nil {
			panic(errutil.Unauthorized("invalid client credentials"))
		}

		c.Set(ProjectContextKey, project)
		c.Next()
	}
}
