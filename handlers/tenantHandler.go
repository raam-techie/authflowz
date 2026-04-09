package handlers

import (
	"new-auth-service/errutil"
	"new-auth-service/payload"
	"new-auth-service/services"

	"github.com/gin-gonic/gin"
)

func CreateTenant(c *gin.Context) {
	var req payload.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	tenant, err := services.CreateTenant(c.Request.Context(), &req)
	if err != nil {
		panic(errutil.Internal(err.Error()))
	}

	c.JSON(201, payload.SuccessResponse{
		StatusCode: 201,
		Message:    "Tenant created successfully",
		Data: []any{map[string]any{
			"tenantId":  tenant.ID,
			"accountId": tenant.AccountID,
			"name":      tenant.Name,
			"email":     tenant.Email,
			"status":    tenant.Status,
			"createdAt": tenant.CreatedAt,
		}},
	})
}

func TenantLogin(c *gin.Context) {
	var req payload.TenantLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	token, err := services.LoginTenant(c.Request.Context(), &req)
	if err != nil {
		msg := err.Error()
		if msg == "invalid credentials" || msg == "tenant account is not active" {
			panic(errutil.Unauthorized(msg))
		}
		panic(errutil.Internal(msg))
	}

	c.JSON(200, payload.SuccessResponse{
		StatusCode: 200,
		Message:    "Login successful",
		Data:       []any{token},
	})
}
