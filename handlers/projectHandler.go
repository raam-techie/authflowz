package handlers

import (
	"net/http"

	"new-auth-service/payload"
	"new-auth-service/services"

	"github.com/gin-gonic/gin"
)

func CreateProject(c *gin.Context) {
	var req payload.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := services.CreateProject(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           project.ID,
		"tenantId":     project.TenantID,
		"clientId":     project.ClientID,
		"name":         project.Name,
		"description":  project.Description,
		"appType":      project.AppType,
		"authType":     project.AuthType,
		"redirectUris": project.RedirectURIs,
		"isActive":     project.IsActive,
		"createdAt":    project.CreatedAt,
	})
}
