package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Action constants
const (
	AuditActionUserLogin           = "user.login"
	AuditActionUserCreate          = "user.create"
	AuditActionUserAttributeUpdate = "user.attributes.update"
	AuditActionOTPSend             = "otp.send"
	AuditActionOTPVerify           = "otp.verify"
	AuditActionTokenRefresh        = "token.refresh"
	AuditActionTenantLogin         = "tenant.login"
)

// ActorType constants
const (
	ActorTypeUser   = "USER"
	ActorTypeTenant = "TENANT"
	ActorTypeSystem = "SYSTEM"
)

// AuditStatus constants
const (
	AuditStatusSuccess = "SUCCESS"
	AuditStatusFailure = "FAILURE"
)

type AuditLog struct {
	ID            string
	TenantID      string
	AppClientID   string
	UserID        string
	Action        string
	ActorType     string
	IPAddress     string
	UserAgent     string
	Status        string
	FailureReason string
	Metadata      map[string]any
	CreatedAt     time.Time
}

func InsertAuditLog(ctx context.Context, entry *AuditLog) error {
	var metaBytes []byte
	if entry.Metadata != nil {
		b, err := json.Marshal(entry.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal audit metadata: %w", err)
		}
		metaBytes = b
	}

	query := `
		INSERT INTO audit_logs
			(tenant_id, app_client_id, user_id, action, actor_type,
			 ip_address, user_agent, status, failure_reason, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := DB.ExecContext(ctx, query,
		nullableString(entry.TenantID),
		nullableString(entry.AppClientID),
		nullableString(entry.UserID),
		entry.Action,
		entry.ActorType,
		nullableString(entry.IPAddress),
		nullableString(entry.UserAgent),
		entry.Status,
		nullableString(entry.FailureReason),
		metaBytes,
	)
	return err
}

type AuditLogFilter struct {
	TenantID    string
	UserID      string
	Action      string
	Status      string
	AppClientID string
	FromDate    time.Time
	ToDate      time.Time
}

func GetAuditLogs(ctx context.Context, filter AuditLogFilter) ([]AuditLog, error) {
	args := []any{filter.TenantID}
	clauses := []string{"tenant_id = $1"}
	i := 2

	if filter.UserID != "" {
		clauses = append(clauses, fmt.Sprintf("user_id = $%d", i))
		args = append(args, filter.UserID)
		i++
	}
	if filter.Action != "" {
		clauses = append(clauses, fmt.Sprintf("action = $%d", i))
		args = append(args, filter.Action)
		i++
	}
	if filter.Status != "" {
		clauses = append(clauses, fmt.Sprintf("status = $%d", i))
		args = append(args, filter.Status)
		i++
	}
	if filter.AppClientID != "" {
		clauses = append(clauses, fmt.Sprintf("app_client_id = $%d", i))
		args = append(args, filter.AppClientID)
		i++
	}
	if !filter.FromDate.IsZero() {
		clauses = append(clauses, fmt.Sprintf("created_at >= $%d", i))
		args = append(args, filter.FromDate)
		i++
	}
	if !filter.ToDate.IsZero() {
		clauses = append(clauses, fmt.Sprintf("created_at <= $%d", i))
		args = append(args, filter.ToDate)
		i++
	}
	_ = i

	query := fmt.Sprintf(`
		SELECT
			id,
			COALESCE(tenant_id::TEXT, ''),
			COALESCE(app_client_id::TEXT, ''),
			COALESCE(user_id::TEXT, ''),
			action,
			actor_type,
			COALESCE(ip_address, ''),
			COALESCE(user_agent, ''),
			status,
			COALESCE(failure_reason, ''),
			metadata,
			created_at
		FROM audit_logs
		WHERE %s
		ORDER BY created_at DESC
	`, strings.Join(clauses, " AND "))

	rows, err := DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var a AuditLog
		var rawMeta []byte
		if err := rows.Scan(
			&a.ID, &a.TenantID, &a.AppClientID, &a.UserID,
			&a.Action, &a.ActorType,
			&a.IPAddress, &a.UserAgent,
			&a.Status, &a.FailureReason,
			&rawMeta, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		if rawMeta != nil {
			_ = json.Unmarshal(rawMeta, &a.Metadata)
		}
		logs = append(logs, a)
	}
	return logs, rows.Err()
}
