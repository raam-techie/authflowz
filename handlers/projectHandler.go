package handlers

import (
	"new-auth-service/errutil"
	"new-auth-service/payload"
	"new-auth-service/services"

	"github.com/gin-gonic/gin"
)

func GetProject(c *gin.Context) {
	id := c.Query("id")
	tenantID := c.Query("tenantId")

	if id == "" && tenantID == "" {
		panic(errutil.BadRequest("either 'id' or 'tenantId' query param is required"))
	}

	if id != "" {
		project, err := services.GetProjectByID(c.Request.Context(), id)
		if err != nil {
			if err.Error() == "project not found" {
				panic(errutil.NotFound(err.Error()))
			}
			panic(errutil.Internal(err.Error()))
		}

		c.JSON(200, payload.SuccessResponse{
			StatusCode: 200,
			Message:    "Project fetched successfully",
			Data: []any{map[string]any{
				"id":           project.ID,
				"tenantId":     project.TenantID,
				"clientId":     project.ClientID,
				"name":         project.Name,
				"description":  project.Description,
				"appType":      project.AppType,
				"authTypes":    project.AuthTypes,
				"redirectUris": project.RedirectURIs,
				"isActive":     project.IsActive,
				"createdAt":    project.CreatedAt,
			}},
		})
		return
	}

	projects, err := services.GetProjectsByTenantID(c.Request.Context(), tenantID)
	if err != nil {
		panic(errutil.Internal(err.Error()))
	}

	data := make([]any, 0, len(projects))
	for _, p := range projects {
		data = append(data, map[string]any{
			"id":           p.ID,
			"tenantId":     p.TenantID,
			"clientId":     p.ClientID,
			"name":         p.Name,
			"description":  p.Description,
			"appType":      p.AppType,
			"authTypes":    p.AuthTypes,
			"redirectUris": p.RedirectURIs,
			"isActive":     p.IsActive,
			"createdAt":    p.CreatedAt,
		})
	}

	c.JSON(200, payload.SuccessResponse{
		StatusCode: 200,
		Message:    "Projects fetched successfully",
		Data:       data,
	})
}

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
