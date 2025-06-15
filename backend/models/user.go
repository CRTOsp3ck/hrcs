package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleNormal UserRole = "normal"
	RoleAdmin  UserRole = "admin"
)

type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Email       string         `json:"email" gorm:"uniqueIndex;not null"`
	Password    string         `json:"-" gorm:"not null"`
	FirstName   string         `json:"first_name" gorm:"not null"`
	LastName    string         `json:"last_name" gorm:"not null"`
	Role        UserRole       `json:"role" gorm:"default:normal"`
	UserGroupID *uint          `json:"user_group_id"`
	UserGroup   *UserGroup     `json:"user_group,omitempty"`
	
	// NEW FIELDS FOR BALANCE TRACKING
	ClaimBalances []UserClaimBalance `json:"claim_balances" gorm:"foreignKey:UserID"`
	
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserGroup struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserClaimBalance model for tracking user balances per claim type
type UserClaimBalance struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"user_id" gorm:"not null"`
	ClaimTypeID     uint      `json:"claim_type_id" gorm:"not null"`
	
	// Balance tracking
	TotalLimit      float64   `json:"total_limit" gorm:"not null"`
	CurrentSpent    float64   `json:"current_spent" gorm:"default:0"`
	RemainingBalance float64  `json:"remaining_balance" gorm:"default:0"`
	
	// Reset tracking
	LastResetDate   time.Time `json:"last_reset_date"`
	ResetPeriod     LimitTimespan `json:"reset_period" gorm:"not null"`
	
	// Relationships
	User            User      `json:"user" gorm:"foreignKey:UserID"`
	ClaimType       ClaimType `json:"claim_type" gorm:"foreignKey:ClaimTypeID"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Helper methods for UserClaimBalance
func (ucb *UserClaimBalance) NeedsReset() bool {
	now := time.Now()
	switch ucb.ResetPeriod {
	case LimitDaily:
		return !isSameDay(ucb.LastResetDate, now)
	case LimitWeekly:
		return !isSameWeek(ucb.LastResetDate, now)
	case LimitMonthly:
		return !isSameMonth(ucb.LastResetDate, now)
	case LimitAnnual:
		return !isSameYear(ucb.LastResetDate, now)
	default:
		return false
	}
}

// Helper functions for date comparison
func isSameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func isSameWeek(date1, date2 time.Time) bool {
	y1, w1 := date1.ISOWeek()
	y2, w2 := date2.ISOWeek()
	return y1 == y2 && w1 == w2
}

func isSameMonth(date1, date2 time.Time) bool {
	y1, m1, _ := date1.Date()
	y2, m2, _ := date2.Date()
	return y1 == y2 && m1 == m2
}

func isSameYear(date1, date2 time.Time) bool {
	return date1.Year() == date2.Year()
}