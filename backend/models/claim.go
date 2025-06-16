package models

import (
	"time"

	"gorm.io/gorm"
)

type ClaimStatus string

const (
	StatusDraft             ClaimStatus = "draft"
	StatusSubmitted         ClaimStatus = "submitted"
	StatusApproved          ClaimStatus = "approved"
	StatusRejected          ClaimStatus = "rejected"
	StatusPaymentInProgress ClaimStatus = "payment-in-progress"
	StatusPaid              ClaimStatus = "paid"
)

type LimitTimespan string

const (
	LimitAnnual  LimitTimespan = "annual"
	LimitMonthly LimitTimespan = "monthly"
	LimitWeekly  LimitTimespan = "weekly"
	LimitDaily   LimitTimespan = "daily"
)

type ClaimType struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`

	// NEW FIELDS FOR CLAIM LIMITS
	LimitAmount   float64       `json:"limit_amount" gorm:"default:0"`
	LimitTimespan LimitTimespan `json:"limit_timespan" gorm:"default:'annual'"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Claim struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	Title       string          `json:"title" gorm:"not null"`
	Description string          `json:"description"`
	Amount      float64         `json:"amount" gorm:"not null"`
	Status      ClaimStatus     `json:"status" gorm:"default:draft"`
	UserID      uint            `json:"user_id" gorm:"not null"`
	User        User            `json:"user"`
	ClaimTypeID uint            `json:"claim_type_id" gorm:"not null"`
	ClaimType   ClaimType       `json:"claim_type"`
	Approvals   []ClaimApproval `json:"approvals,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `json:"-" gorm:"index"`
}

type ApprovalLevel struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Level       int       `json:"level" gorm:"not null"`
	UserGroupID uint      `json:"user_group_id" gorm:"not null"`
	UserGroup   UserGroup `json:"user_group"`
	ApproverID  uint      `json:"approver_id" gorm:"not null"`
	Approver    User      `json:"approver"`
	// Status permissions - what statuses this level can set
	CanDraft                bool           `json:"can_draft" gorm:"default:false"`
	CanSubmit               bool           `json:"can_submit" gorm:"default:false"`
	CanApprove              bool           `json:"can_approve" gorm:"default:true"`
	CanReject               bool           `json:"can_reject" gorm:"default:true"`
	CanSetPaymentInProgress bool           `json:"can_set_payment_in_progress" gorm:"default:false"`
	CanSetPaid              bool           `json:"can_set_paid" gorm:"default:false"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	DeletedAt               gorm.DeletedAt `json:"-" gorm:"index"`
}

type ClaimApproval struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	ClaimID         uint           `json:"claim_id" gorm:"not null"`
	Claim           Claim          `json:"claim"`
	ApprovalLevelID uint           `json:"approval_level_id" gorm:"not null"`
	ApprovalLevel   ApprovalLevel  `json:"approval_level"`
	ApproverID      uint           `json:"approver_id" gorm:"not null"`
	Approver        User           `json:"approver"`
	Status          ClaimStatus    `json:"status" gorm:"not null"`
	Comments        string         `json:"comments"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserGroupClaimType association model for user group claim type permissions
type UserGroupClaimType struct {
	ID                uint     `json:"id" gorm:"primaryKey"`
	UserGroupID       uint     `json:"user_group_id" gorm:"not null"`
	ClaimTypeID       uint     `json:"claim_type_id" gorm:"not null"`
	IsAllowed         bool     `json:"is_allowed" gorm:"default:true"`
	CustomLimitAmount *float64 `json:"custom_limit_amount"`

	UserGroup UserGroup `json:"user_group" gorm:"foreignKey:UserGroupID"`
	ClaimType ClaimType `json:"claim_type" gorm:"foreignKey:ClaimTypeID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserClaimType model for individual user overrides
type UserClaimType struct {
	ID                uint     `json:"id" gorm:"primaryKey"`
	UserID            uint     `json:"user_id" gorm:"not null"`
	ClaimTypeID       uint     `json:"claim_type_id" gorm:"not null"`
	IsAllowed         bool     `json:"is_allowed" gorm:"default:true"`
	CustomLimitAmount *float64 `json:"custom_limit_amount"`

	User      User      `json:"user" gorm:"foreignKey:UserID"`
	ClaimType ClaimType `json:"claim_type" gorm:"foreignKey:ClaimTypeID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
