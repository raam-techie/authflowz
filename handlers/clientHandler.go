package handlers

import (
	"net/http"

	"new-auth-service/payload"
	"new-auth-service/services"

	"github.com/gin-gonic/gin"
)

func CreateClient(c *gin.Context) {
	var req payload.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := services.CreateClient(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":            client.ID,
		"tenant_id":     client.TenantID,
		"client_id":     client.ClientID,
		"name":          client.Name,
		"description":   client.Description,
		"app_type":      client.AppType,
		"auth_type":     client.AuthType,
		"redirect_uris": client.RedirectURIs,
		"is_active":     client.IsActive,
		"created_at":    client.CreatedAt,
	})
}
