package handlers

import (
	"new-auth-service/errutil"
	"new-auth-service/payload"
	"new-auth-service/services"

	"github.com/gin-gonic/gin"
)

func CreateProject(c *gin.Context) {
	var req payload.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	project, clientSecret, err := services.CreateProject(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid tenant id" {
			panic(errutil.NotFound(err.Error()))
		}
		panic(errutil.Internal(err.Error()))
	}

	c.JSON(201, payload.SuccessResponse{
		StatusCode: 201,
		Message:    "Project created successfully",
		Data: []any{map[string]any{
			"id":           project.ID,
			"tenantId":     project.TenantID,
			"clientId":     project.ClientID,
			"clientSecret": clientSecret,
			"name":         project.Name,
			"description":  project.Description,
			"appType":      project.AppType,
			"authTypes":    project.AuthTypes,
			"redirectUris": project.RedirectURIs,
			"isActive":     project.IsActive,
			"createdAt":    project.CreatedAt,
		}},
	})
}
