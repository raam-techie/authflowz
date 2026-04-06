package handlers

import (
	"net/http"

	"new-auth-service/payload"
	"new-auth-service/services"

	"github.com/gin-gonic/gin"
)

func CreateTenant(c *gin.Context) {
	var req payload.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, err := services.CreateTenant(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         tenant.ID,
		"account_id": tenant.AccountID,
		"name":       tenant.Name,
		"email":      tenant.Email,
		"status":     tenant.Status,
		"created_at": tenant.CreatedAt,
	})
}
