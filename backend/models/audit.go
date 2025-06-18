package models

import (
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	User       User      `json:"user"`
	Action     string    `json:"action" gorm:"not null"`
	EntityType string    `json:"entity_type" gorm:"not null"`
	EntityID   uint      `json:"entity_id"`
	OldValues  string    `json:"old_values" gorm:"type:jsonb"`
	NewValues  string    `json:"new_values" gorm:"type:jsonb"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}

// AuditLogRequest represents the request for creating audit logs
type AuditLogRequest struct {
	UserID     uint                   `json:"user_id"`
	Action     string                 `json:"action"`
	EntityType string                 `json:"entity_type"`
	EntityID   uint                   `json:"entity_id,omitempty"`
	OldValues  map[string]interface{} `json:"old_values,omitempty"`
	NewValues  map[string]interface{} `json:"new_values,omitempty"`
	IPAddress  string                 `json:"ip_address,omitempty"`
	UserAgent  string                 `json:"user_agent,omitempty"`
}

// AuditLogResponse represents the response for audit log queries
type AuditLogResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	UserName     string    `json:"user_name"`
	UserEmail    string    `json:"user_email"`
	Action       string    `json:"action"`
	EntityType   string    `json:"entity_type"`
	EntityID     uint      `json:"entity_id"`
	OldValues    string    `json:"old_values"`
	NewValues    string    `json:"new_values"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	CreatedAt    time.Time `json:"created_at"`
	Description  string    `json:"description"`
	ActionType   string    `json:"action_type"`
	Severity     string    `json:"severity"`
}

// AuditLogFilter represents filters for audit log queries
type AuditLogFilter struct {
	UserID     *uint     `json:"user_id,omitempty"`
	Action     string    `json:"action,omitempty"`
	EntityType string    `json:"entity_type,omitempty"`
	StartDate  time.Time `json:"start_date,omitempty"`
	EndDate    time.Time `json:"end_date,omitempty"`
	Search     string    `json:"search,omitempty"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
}

// Common audit actions
const (
	ActionCreate   = "create"
	ActionUpdate   = "update"
	ActionDelete   = "delete"
	ActionApprove  = "approve"
	ActionReject   = "reject"
	ActionSubmit   = "submit"
	ActionLogin    = "login"
	ActionLogout   = "logout"
	ActionView     = "view"
	ActionExport   = "export"
	ActionImport   = "import"
)

// Common entity types
const (
	EntityTypeClaim         = "claim"
	EntityTypeUser          = "user"
	EntityTypeUserGroup     = "user_group"
	EntityTypeClaimType     = "claim_type"
	EntityTypeApprovalLevel = "approval_level"
	EntityTypeClaimApproval = "claim_approval"
	EntityTypeSystem        = "system"
)

// Helper methods for AuditLog
func (al *AuditLog) GetDescription() string {
	switch al.Action {
	case ActionCreate:
		return "Created " + al.EntityType
	case ActionUpdate:
		return "Updated " + al.EntityType
	case ActionDelete:
		return "Deleted " + al.EntityType
	case ActionApprove:
		return "Approved " + al.EntityType
	case ActionReject:
		return "Rejected " + al.EntityType
	case ActionSubmit:
		return "Submitted " + al.EntityType
	case ActionLogin:
		return "User logged in"
	case ActionLogout:
		return "User logged out"
	case ActionView:
		return "Viewed " + al.EntityType
	case ActionExport:
		return "Exported " + al.EntityType + " data"
	case ActionImport:
		return "Imported " + al.EntityType + " data"
	default:
		return al.Action + " " + al.EntityType
	}
}

func (al *AuditLog) GetActionType() string {
	switch al.Action {
	case ActionCreate, ActionImport:
		return "create"
	case ActionUpdate, ActionApprove, ActionReject, ActionSubmit:
		return "update"
	case ActionDelete:
		return "delete"
	case ActionLogin, ActionLogout:
		return "auth"
	case ActionView, ActionExport:
		return "read"
	default:
		return "other"
	}
}

func (al *AuditLog) GetSeverity() string {
	switch al.Action {
	case ActionDelete:
		return "high"
	case ActionApprove, ActionReject, ActionCreate, ActionUpdate:
		return "medium"
	case ActionLogin, ActionLogout, ActionView, ActionExport:
		return "low"
	default:
		return "medium"
	}
}

// BeforeCreate hook for GORM
func (al *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if al.CreatedAt.IsZero() {
		al.CreatedAt = time.Now()
	}
	return nil
}