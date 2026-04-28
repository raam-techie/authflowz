package payload

type AuditLogQueryParams struct {
	TenantID    string `form:"tenant_id" binding:"required"`
	UserID      string `form:"user_id"`
	Action      string `form:"action"`
	Status      string `form:"status"`
	AppClientID string `form:"app_client_id"`
	FromDate    string `form:"from_date"`
	ToDate      string `form:"to_date"`
}
