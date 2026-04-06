package handlers

import (
	"net/http"

	"new-auth-service/payload"
	"new-auth-service/services"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var req payload.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"tenant_id":  user.TenantID,
		"project_id": user.ProjectID,
		"uid":        user.UID,
		"name":       user.Name,
		"email":      user.Email,
		"phone":      user.Phone,
		"department": user.Department,
		"role":       user.Role,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
	})
}
