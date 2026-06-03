package handlers

import (
	"time"

	"new-auth-service/internal/db"
	"new-auth-service/internal/errutil"
	"new-auth-service/internal/payload"

	"github.com/gin-gonic/gin"
)

func GetAuditLogs(c *gin.Context) {
	var params payload.AuditLogQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	filter := db.AuditLogFilter{
		TenantID:    params.TenantID,
		UserID:      params.UserID,
		Action:      params.Action,
		Status:      params.Status,
		AppClientID: params.AppClientID,
	}

	if params.FromDate != "" {
		t, err := time.Parse(time.RFC3339, params.FromDate)
		if err != nil {
			panic(errutil.BadRequest("invalid from_date format, use ISO8601 (e.g. 2025-01-01T00:00:00Z)"))
		}
		filter.FromDate = t
	}

	if params.ToDate != "" {
		t, err := time.Parse(time.RFC3339, params.ToDate)
		if err != nil {
			panic(errutil.BadRequest("invalid to_date format, use ISO8601 (e.g. 2025-01-01T00:00:00Z)"))
		}
		filter.ToDate = t
	}

	logs, err := db.GetAuditLogs(c.Request.Context(), filter)
	if err != nil {
		panic(errutil.Internal(err.Error()))
	}

	data := make([]any, 0, len(logs))
	for _, l := range logs {
		data = append(data, map[string]any{
			"id":            l.ID,
			"tenantId":      l.TenantID,
			"appClientId":   l.AppClientID,
			"userId":        l.UserID,
			"action":        l.Action,
			"actorType":     l.ActorType,
			"ipAddress":     l.IPAddress,
			"userAgent":     l.UserAgent,
			"status":        l.Status,
			"failureReason": l.FailureReason,
			"metadata":      l.Metadata,
			"createdAt":     l.CreatedAt,
		})
	}

	c.JSON(200, payload.SuccessResponse{
		StatusCode: 200,
		Message:    "Audit logs fetched successfully",
		Data:       data,
	})
}
