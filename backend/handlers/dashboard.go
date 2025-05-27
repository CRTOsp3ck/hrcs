package handlers

import (
	"hrcs/backend/middleware"
	"hrcs/backend/models"
	"hrcs/backend/utils"
	"net/http"

	"gorm.io/gorm"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

type DashboardStats struct {
	TotalClaims    int64              `json:"totalClaims"`
	PendingClaims  int64              `json:"pendingClaims"`
	ApprovedClaims int64              `json:"approvedClaims"`
	RejectedClaims int64              `json:"rejectedClaims"`
	TotalAmount    float64            `json:"totalAmount"`
	ApprovedAmount float64            `json:"approvedAmount"`
	RecentClaims   []models.Claim     `json:"recentClaims"`
	ClaimsByStatus []ClaimStatusCount `json:"claimsByStatus"`
	ClaimsByType   []ClaimTypeStats   `json:"claimsByType"`
	TotalUsers     int64              `json:"totalUsers"` // Admin only
}

type ClaimStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type ClaimTypeStats struct {
	Type   string  `json:"type"`
	Count  int64   `json:"count"`
	Amount float64 `json:"amount"`
}

func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserContextKey).(*models.User)
	if user == nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not found in context")
		return
	}
	userID := user.ID

	var stats DashboardStats

	// Get total claims for user
	h.db.Model(&models.Claim{}).Where("user_id = ?", userID).Count(&stats.TotalClaims)

	// Get pending claims (submitted status)
	h.db.Model(&models.Claim{}).Where("user_id = ? AND status = ?", userID, "submitted").Count(&stats.PendingClaims)

	// Get approved claims
	h.db.Model(&models.Claim{}).Where("user_id = ? AND status = ?", userID, "approved").Count(&stats.ApprovedClaims)

	// Get rejected claims
	h.db.Model(&models.Claim{}).Where("user_id = ? AND status = ?", userID, "rejected").Count(&stats.RejectedClaims)

	// Get total amount
	h.db.Model(&models.Claim{}).Where("user_id = ?", userID).Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalAmount)

	// Get approved amount
	h.db.Model(&models.Claim{}).Where("user_id = ? AND status = ?", userID, "approved").Select("COALESCE(SUM(amount), 0)").Scan(&stats.ApprovedAmount)

	// Get recent claims
	var recentClaims []models.Claim
	h.db.Preload("ClaimType").Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Limit(5).Find(&recentClaims)
	stats.RecentClaims = recentClaims

	// Get claims by status
	var claimsByStatus []ClaimStatusCount
	h.db.Model(&models.Claim{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Scan(&claimsByStatus)
	stats.ClaimsByStatus = claimsByStatus

	// Get claims by type
	var claimsByType []ClaimTypeStats
	h.db.Table("claims").
		Select("claim_types.name as type, COUNT(claims.id) as count, COALESCE(SUM(claims.amount), 0) as amount").
		Joins("JOIN claim_types ON claim_types.id = claims.claim_type_id").
		Where("claims.user_id = ?", userID).
		Group("claim_types.id, claim_types.name").
		Scan(&claimsByType)
	stats.ClaimsByType = claimsByType

	utils.WriteSuccess(w, stats, "Dashboard stats retrieved successfully")
}

func (h *DashboardHandler) GetAdminStats(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	user := r.Context().Value(middleware.UserContextKey).(*models.User)
	if user == nil || user.Role != models.RoleAdmin {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var stats DashboardStats

	// Get total claims (all users)
	h.db.Model(&models.Claim{}).Count(&stats.TotalClaims)

	// Get pending claims
	h.db.Model(&models.Claim{}).Where("status = ?", "submitted").Count(&stats.PendingClaims)

	// Get approved claims
	h.db.Model(&models.Claim{}).Where("status = ?", "approved").Count(&stats.ApprovedClaims)

	// Get rejected claims
	h.db.Model(&models.Claim{}).Where("status = ?", "rejected").Count(&stats.RejectedClaims)

	// Get total amount
	h.db.Model(&models.Claim{}).Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalAmount)

	// Get approved amount
	h.db.Model(&models.Claim{}).Where("status = ?", "approved").Select("COALESCE(SUM(amount), 0)").Scan(&stats.ApprovedAmount)

	// Get recent claims
	var recentClaims []models.Claim
	h.db.Preload("ClaimType").Preload("User").Order("created_at DESC").Limit(10).Find(&recentClaims)
	stats.RecentClaims = recentClaims

	// Get claims by status
	var claimsByStatus []ClaimStatusCount
	h.db.Model(&models.Claim{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&claimsByStatus)
	stats.ClaimsByStatus = claimsByStatus

	// Get claims by type
	var claimsByType []ClaimTypeStats
	h.db.Table("claims").
		Select("claim_types.name as type, COUNT(claims.id) as count, COALESCE(SUM(claims.amount), 0) as amount").
		Joins("JOIN claim_types ON claim_types.id = claims.claim_type_id").
		Group("claim_types.id, claim_types.name").
		Scan(&claimsByType)
	stats.ClaimsByType = claimsByType

	// Get total users count
	h.db.Model(&models.User{}).Count(&stats.TotalUsers)

	utils.WriteSuccess(w, stats, "Admin dashboard stats retrieved successfully")
}
