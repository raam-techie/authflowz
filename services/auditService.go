package services

import (
	"context"
	"log"

	"new-auth-service/db"
)

type AuditEntry struct {
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
}

// LogAction writes an audit log entry in a background goroutine.
// It never blocks the caller and never propagates errors.
func LogAction(entry AuditEntry) {
	go func() {
		logEntry := &db.AuditLog{
			TenantID:      entry.TenantID,
			AppClientID:   entry.AppClientID,
			UserID:        entry.UserID,
			Action:        entry.Action,
			ActorType:     entry.ActorType,
			IPAddress:     entry.IPAddress,
			UserAgent:     entry.UserAgent,
			Status:        entry.Status,
			FailureReason: entry.FailureReason,
			Metadata:      entry.Metadata,
		}
		if err := db.InsertAuditLog(context.Background(), logEntry); err != nil {
			log.Printf("[audit] failed to insert log action=%s: %v", entry.Action, err)
		}
	}()
}
