package handlers

import (
	"errors"

	"new-auth-service/errutil"
	"new-auth-service/payload"
	"new-auth-service/services"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var req payload.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	user, err := services.CreateUser(c.Request.Context(), &req)
	if err != nil {
		panic(errutil.Internal(err.Error()))
	}

	c.JSON(201, payload.SuccessResponse{
		StatusCode: 201,
		Message:    "User created successfully",
		Data: []any{map[string]any{
			"id":         user.ID,
			"tenantId":   user.TenantID,
			"projectId":  user.ProjectID,
			"userId":     user.UID,
			"name":       user.Name,
			"email":      user.Email,
			"phone":      user.Phone,
			"department": user.Department,
			"role":       user.Role,
			"isActive":   user.IsActive,
			"createdAt":  user.CreatedAt,
		}},
	})
}

func UserLogin(c *gin.Context) {
	var req payload.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	token, err := services.LoginUser(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrNotImplemented) {
			panic(errutil.NotImplemented(err.Error()))
		}
		msg := err.Error()
		if msg == "invalid credentials" || msg == "user account is not active" || msg == "project is not active" {
			panic(errutil.Unauthorized(msg))
		}
		if msg == "project not found" {
			panic(errutil.NotFound(msg))
		}
		panic(errutil.Internal(msg))
	}

	c.JSON(200, payload.SuccessResponse{
		StatusCode: 200,
		Message:    "Login successful",
		Data:       []any{token},
	})
}
