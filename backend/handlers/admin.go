package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"hrcs/backend/models"
	"hrcs/backend/utils"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type AdminHandler struct {
	DB *gorm.DB
}

type CreateClaimTypeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateUserGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateApprovalLevelRequest struct {
	Level                   int  `json:"level"`
	UserGroupID             uint `json:"user_group_id"`
	ApproverID              uint `json:"approver_id"`
	CanDraft                bool `json:"can_draft"`
	CanSubmit               bool `json:"can_submit"`
	CanApprove              bool `json:"can_approve"`
	CanReject               bool `json:"can_reject"`
	CanSetPaymentInProgress bool `json:"can_set_payment_in_progress"`
	CanSetPaid              bool `json:"can_set_paid"`
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

func (h *AdminHandler) CreateClaimType(w http.ResponseWriter, r *http.Request) {
	var req CreateClaimTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	claimType := models.ClaimType{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.DB.Create(&claimType).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create claim type")
		return
	}

	utils.WriteSuccess(w, claimType, "Claim type created successfully")
}

func (h *AdminHandler) GetClaimTypes(w http.ResponseWriter, r *http.Request) {
	var claimTypes []models.ClaimType
	if err := h.DB.Find(&claimTypes).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve claim types")
		return
	}

	utils.WriteSuccess(w, claimTypes)
}

func (h *AdminHandler) UpdateClaimType(w http.ResponseWriter, r *http.Request) {
	claimTypeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid claim type ID")
		return
	}

	var claimType models.ClaimType
	if err := h.DB.First(&claimType, claimTypeID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Claim type not found")
		return
	}

	var req CreateClaimTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	claimType.Name = req.Name
	claimType.Description = req.Description

	if err := h.DB.Save(&claimType).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update claim type")
		return
	}

	utils.WriteSuccess(w, claimType, "Claim type updated successfully")
}

func (h *AdminHandler) DeleteClaimType(w http.ResponseWriter, r *http.Request) {
	claimTypeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid claim type ID")
		return
	}

	if err := h.DB.Delete(&models.ClaimType{}, claimTypeID).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete claim type")
		return
	}

	utils.WriteSuccess(w, nil, "Claim type deleted successfully")
}

func (h *AdminHandler) CreateUserGroup(w http.ResponseWriter, r *http.Request) {
	var req CreateUserGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userGroup := models.UserGroup{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.DB.Create(&userGroup).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create user group")
		return
	}

	utils.WriteSuccess(w, userGroup, "User group created successfully")
}

func (h *AdminHandler) GetUserGroups(w http.ResponseWriter, r *http.Request) {
	var userGroups []models.UserGroup
	if err := h.DB.Find(&userGroups).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve user groups")
		return
	}

	utils.WriteSuccess(w, userGroups)
}

func (h *AdminHandler) UpdateUserGroup(w http.ResponseWriter, r *http.Request) {
	userGroupID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user group ID")
		return
	}

	var userGroup models.UserGroup
	if err := h.DB.First(&userGroup, userGroupID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "User group not found")
		return
	}

	var req CreateUserGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userGroup.Name = req.Name
	userGroup.Description = req.Description

	if err := h.DB.Save(&userGroup).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update user group")
		return
	}

	utils.WriteSuccess(w, userGroup, "User group updated successfully")
}

func (h *AdminHandler) DeleteUserGroup(w http.ResponseWriter, r *http.Request) {
	userGroupID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user group ID")
		return
	}

	if err := h.DB.Delete(&models.UserGroup{}, userGroupID).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete user group")
		return
	}

	utils.WriteSuccess(w, nil, "User group deleted successfully")
}

func (h *AdminHandler) CreateApprovalLevel(w http.ResponseWriter, r *http.Request) {
	var req CreateApprovalLevelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	approvalLevel := models.ApprovalLevel{
		Level:                   req.Level,
		UserGroupID:             req.UserGroupID,
		ApproverID:              req.ApproverID,
		CanDraft:                req.CanDraft,
		CanSubmit:               req.CanSubmit,
		CanApprove:              req.CanApprove,
		CanReject:               req.CanReject,
		CanSetPaymentInProgress: req.CanSetPaymentInProgress,
		CanSetPaid:              req.CanSetPaid,
	}

	if err := h.DB.Create(&approvalLevel).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create approval level")
		return
	}

	if err := h.DB.Preload("UserGroup").Preload("Approver").First(&approvalLevel, approvalLevel.ID).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve approval level")
		return
	}

	utils.WriteSuccess(w, approvalLevel, "Approval level created successfully")
}

func (h *AdminHandler) GetApprovalLevels(w http.ResponseWriter, r *http.Request) {
	var approvalLevels []models.ApprovalLevel
	if err := h.DB.Preload("UserGroup").Preload("Approver").Find(&approvalLevels).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve approval levels")
		return
	}

	utils.WriteSuccess(w, approvalLevels)
}

func (h *AdminHandler) UpdateApprovalLevel(w http.ResponseWriter, r *http.Request) {
	approvalLevelID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid approval level ID")
		return
	}

	var approvalLevel models.ApprovalLevel
	if err := h.DB.First(&approvalLevel, approvalLevelID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Approval level not found")
		return
	}

	var req CreateApprovalLevelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	approvalLevel.Level = req.Level
	approvalLevel.UserGroupID = req.UserGroupID
	approvalLevel.ApproverID = req.ApproverID
	approvalLevel.CanDraft = req.CanDraft
	approvalLevel.CanSubmit = req.CanSubmit
	approvalLevel.CanApprove = req.CanApprove
	approvalLevel.CanReject = req.CanReject
	approvalLevel.CanSetPaymentInProgress = req.CanSetPaymentInProgress
	approvalLevel.CanSetPaid = req.CanSetPaid

	if err := h.DB.Save(&approvalLevel).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update approval level")
		return
	}

	if err := h.DB.Preload("UserGroup").Preload("Approver").First(&approvalLevel, approvalLevel.ID).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve updated approval level")
		return
	}

	utils.WriteSuccess(w, approvalLevel, "Approval level updated successfully")
}

func (h *AdminHandler) DeleteApprovalLevel(w http.ResponseWriter, r *http.Request) {
	approvalLevelID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid approval level ID")
		return
	}

	if err := h.DB.Delete(&models.ApprovalLevel{}, approvalLevelID).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete approval level")
		return
	}

	utils.WriteSuccess(w, nil, "Approval level deleted successfully")
}

// NEW ENDPOINTS FOR CLAIM TYPE LIMITS AND PERMISSIONS

type UpdateClaimTypeLimitsRequest struct {
	LimitAmount   float64          `json:"limit_amount"`
	LimitTimespan models.LimitTimespan `json:"limit_timespan"`
}

func (h *AdminHandler) UpdateClaimTypeWithLimits(w http.ResponseWriter, r *http.Request) {
	claimTypeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid claim type ID")
		return
	}

	var claimType models.ClaimType
	if err := h.DB.First(&claimType, claimTypeID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Claim type not found")
		return
	}

	var req UpdateClaimTypeLimitsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	claimType.LimitAmount = req.LimitAmount
	claimType.LimitTimespan = req.LimitTimespan

	if err := h.DB.Save(&claimType).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update claim type limits")
		return
	}

	utils.WriteSuccess(w, claimType, "Claim type limits updated successfully")
}

type SetUserGroupClaimPermissionsRequest struct {
	Permissions []struct {
		ClaimTypeID       uint     `json:"claim_type_id"`
		IsAllowed         bool     `json:"is_allowed"`
		CustomLimitAmount *float64 `json:"custom_limit_amount"`
	} `json:"permissions"`
}

func (h *AdminHandler) SetUserGroupClaimPermissions(w http.ResponseWriter, r *http.Request) {
	userGroupID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user group ID")
		return
	}

	var req SetUserGroupClaimPermissionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Begin transaction
	tx := h.DB.Begin()
	defer tx.Rollback()

	// Delete existing permissions for this group
	if err := tx.Where("user_group_id = ?", userGroupID).Delete(&models.UserGroupClaimType{}).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to clear existing permissions")
		return
	}

	// Create new permissions
	for _, perm := range req.Permissions {
		permission := models.UserGroupClaimType{
			UserGroupID:       uint(userGroupID),
			ClaimTypeID:       perm.ClaimTypeID,
			IsAllowed:         perm.IsAllowed,
			CustomLimitAmount: perm.CustomLimitAmount,
		}

		if err := tx.Create(&permission).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to create permission")
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to save permissions")
		return
	}

	utils.WriteSuccess(w, nil, "User group claim permissions updated successfully")
}

type SetUserClaimOverridesRequest struct {
	Overrides []struct {
		ClaimTypeID       uint     `json:"claim_type_id"`
		IsAllowed         bool     `json:"is_allowed"`
		CustomLimitAmount *float64 `json:"custom_limit_amount"`
	} `json:"overrides"`
}

func (h *AdminHandler) SetUserClaimOverrides(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req SetUserClaimOverridesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Begin transaction
	tx := h.DB.Begin()
	defer tx.Rollback()

	// Delete existing overrides for this user
	if err := tx.Where("user_id = ?", userID).Delete(&models.UserClaimType{}).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to clear existing overrides")
		return
	}

	// Create new overrides
	for _, override := range req.Overrides {
		userOverride := models.UserClaimType{
			UserID:            uint(userID),
			ClaimTypeID:       override.ClaimTypeID,
			IsAllowed:         override.IsAllowed,
			CustomLimitAmount: override.CustomLimitAmount,
		}

		if err := tx.Create(&userOverride).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to create override")
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to save overrides")
		return
	}

	utils.WriteSuccess(w, nil, "User claim overrides updated successfully")
}

// NEW DETAIL VIEW ENDPOINTS FOR PHASE 2

// GetUserDetails gets detailed information about a specific user
func (h *AdminHandler) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := h.DB.Preload("UserGroup").First(&user, userID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	// Get user's claim balances
	var balances []models.UserClaimBalance
	h.DB.Where("user_id = ?", userID).Preload("ClaimType").Find(&balances)

	// Get recent claims
	var claims []models.Claim
	h.DB.Where("user_id = ?", userID).Preload("ClaimType").Order("created_at DESC").Limit(10).Find(&claims)

	// Get user-specific claim permissions
	var permissions []models.UserClaimType
	h.DB.Where("user_id = ?", userID).Preload("ClaimType").Find(&permissions)

	response := map[string]interface{}{
		"user":        user,
		"balances":    balances,
		"claims":      claims,
		"permissions": permissions,
	}

	utils.WriteSuccess(w, response)
}

// GetUserGroupDetails gets detailed information about a specific user group
func (h *AdminHandler) GetUserGroupDetails(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	var group models.UserGroup
	if err := h.DB.First(&group, groupID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "User group not found")
		return
	}

	// Get group members
	var members []models.User
	h.DB.Where("user_group_id = ?", groupID).Find(&members)

	// Get group claim permissions
	var permissions []models.UserGroupClaimType
	h.DB.Where("user_group_id = ?", groupID).Preload("ClaimType").Find(&permissions)

	// Get approval levels for this group
	var approvalLevels []models.ApprovalLevel
	h.DB.Where("user_group_id = ?", groupID).Preload("Approver").Find(&approvalLevels)

	response := map[string]interface{}{
		"group":           group,
		"members":         members,
		"permissions":     permissions,
		"approval_levels": approvalLevels,
	}

	utils.WriteSuccess(w, response)
}

// GetClaimTypeDetails gets detailed information about a specific claim type
func (h *AdminHandler) GetClaimTypeDetails(w http.ResponseWriter, r *http.Request) {
	claimTypeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid claim type ID")
		return
	}

	var claimType models.ClaimType
	if err := h.DB.First(&claimType, claimTypeID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Claim type not found")
		return
	}

	// Get statistics for this claim type
	type ClaimTypeStats struct {
		TotalClaims    int64   `json:"total_claims"`
		ApprovedClaims int64   `json:"approved_claims"`
		TotalAmount    float64 `json:"total_amount"`
		AverageAmount  float64 `json:"average_amount"`
	}

	var stats ClaimTypeStats
	h.DB.Model(&models.Claim{}).Where("claim_type_id = ?", claimTypeID).Count(&stats.TotalClaims)
	h.DB.Model(&models.Claim{}).Where("claim_type_id = ? AND status = ?", claimTypeID, models.StatusApproved).Count(&stats.ApprovedClaims)
	h.DB.Model(&models.Claim{}).Where("claim_type_id = ?", claimTypeID).Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalAmount)
	
	if stats.TotalClaims > 0 {
		stats.AverageAmount = stats.TotalAmount / float64(stats.TotalClaims)
	}

	// Get user group permissions for this claim type
	var groupPermissions []models.UserGroupClaimType
	h.DB.Where("claim_type_id = ?", claimTypeID).Preload("UserGroup").Find(&groupPermissions)

	// Get recent claims for this type
	var recentClaims []models.Claim
	h.DB.Where("claim_type_id = ?", claimTypeID).Preload("User").Order("created_at DESC").Limit(10).Find(&recentClaims)

	response := map[string]interface{}{
		"claim_type":        claimType,
		"stats":            stats,
		"group_permissions": groupPermissions,
		"recent_claims":     recentClaims,
	}

	utils.WriteSuccess(w, response)
}