package handlers

import (
	"net/http"
	"new-auth-service/internal/db"
	"new-auth-service/internal/errutil"
	"new-auth-service/internal/payload"
	"new-auth-service/internal/services"

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

	c.JSON(http.StatusCreated, payload.SuccessResponse{
		StatusCode: http.StatusCreated,
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

	tokens, err := services.LoginTenant(c.Request.Context(), &req)

	auditStatus := db.AuditStatusSuccess
	auditReason := ""
	if err != nil {
		auditStatus = db.AuditStatusFailure
		auditReason = err.Error()
	}
	services.LogAction(services.AuditEntry{
		Action:        db.AuditActionTenantLogin,
		ActorType:     db.ActorTypeTenant,
		IPAddress:     c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Status:        auditStatus,
		FailureReason: auditReason,
		Metadata:      map[string]any{"accountId": req.AccountID},
	})

	if err != nil {
		msg := err.Error()
		if msg == "invalid credentials" || msg == "tenant account is not active" {
			panic(errutil.Unauthorized(msg))
		}
		panic(errutil.Internal(msg))
	}

	c.JSON(http.StatusOK, payload.SuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "Login successful",
		Data:       []any{tokens},
	})
}
