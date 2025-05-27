package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"hrcs/backend/models"
	"hrcs/backend/utils"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type AdminEnhancedHandler struct {
	DB *gorm.DB
}

func NewAdminEnhancedHandler(db *gorm.DB) *AdminEnhancedHandler {
	return &AdminEnhancedHandler{DB: db}
}

// Admin Claims Management
func (h *AdminEnhancedHandler) GetAllClaims(w http.ResponseWriter, r *http.Request) {
	var claims []models.Claim
	query := h.DB.Preload("User").Preload("ClaimType").Order("created_at DESC")

	// Add filters if needed
	status := r.URL.Query().Get("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&claims).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve claims")
		return
	}

	// Enhanced claim response with additional fields
	type EnhancedClaim struct {
		models.Claim
		Employee         string `json:"employee"`
		Department       string `json:"department"`
		Type             string `json:"type"`
		ApprovalsReceived int   `json:"approvalsReceived"`
		ApprovalsRequired int   `json:"approvalsRequired"`
		SubmittedDate    string `json:"submittedDate"`
	}

	var enhancedClaims []EnhancedClaim
	for _, claim := range claims {
		enhanced := EnhancedClaim{
			Claim:            claim,
			Employee:         claim.User.FirstName + " " + claim.User.LastName,
			Department:       "IT", // TODO: Add department field to User model
			Type:             claim.ClaimType.Name,
			ApprovalsReceived: 1, // TODO: Implement actual approval tracking
			ApprovalsRequired: 3, // TODO: Based on approval levels
			SubmittedDate:    claim.CreatedAt.Format("2006-01-02"),
		}
		enhancedClaims = append(enhancedClaims, enhanced)
	}

	utils.WriteSuccess(w, enhancedClaims)
}

func (h *AdminEnhancedHandler) AdminApproveClaim(w http.ResponseWriter, r *http.Request) {
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	var req struct {
		Comments string `json:"comments"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var claim models.Claim
	if err := h.DB.First(&claim, claimID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Claim not found")
		return
	}

	claim.Status = models.StatusApproved
	if err := h.DB.Save(&claim).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to approve claim")
		return
	}

	utils.WriteSuccess(w, claim, "Claim approved successfully")
}

func (h *AdminEnhancedHandler) AdminRejectClaim(w http.ResponseWriter, r *http.Request) {
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	var req struct {
		Comments string `json:"comments"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var claim models.Claim
	if err := h.DB.First(&claim, claimID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Claim not found")
		return
	}

	claim.Status = models.StatusRejected
	if err := h.DB.Save(&claim).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to reject claim")
		return
	}

	utils.WriteSuccess(w, claim, "Claim rejected successfully")
}

// Admin Users Management
func (h *AdminEnhancedHandler) GetAdminUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := h.DB.Preload("UserGroup").Find(&users).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	// Enhanced user response
	type EnhancedUser struct {
		ID         uint               `json:"id"`
		Email      string             `json:"email"`
		Name       string             `json:"name"`
		FirstName  string             `json:"first_name"`
		LastName   string             `json:"last_name"`
		Role       models.UserRole    `json:"role"`
		Department string             `json:"department"`
		Groups     []models.UserGroup `json:"groups"`
		Status     string             `json:"status"`
		CreatedAt  time.Time          `json:"createdAt"`
	}

	var enhancedUsers []EnhancedUser
	for _, user := range users {
		groups := []models.UserGroup{}
		if user.UserGroup != nil {
			groups = append(groups, *user.UserGroup)
		}
		
		enhanced := EnhancedUser{
			ID:         user.ID,
			Email:      user.Email,
			Name:       user.FirstName + " " + user.LastName,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Role:       user.Role,
			Department: "IT", // TODO: Add department field to User model
			Groups:     groups,
			Status:     "active", // TODO: Implement user status
			CreatedAt:  user.CreatedAt,
		}
		enhancedUsers = append(enhancedUsers, enhanced)
	}

	utils.WriteSuccess(w, enhancedUsers)
}

func (h *AdminEnhancedHandler) CreateAdminUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name       string   `json:"name"`
		Email      string   `json:"email"`
		Password   string   `json:"password"`
		Role       string   `json:"role"`
		Department string   `json:"department"`
		Groups     []uint   `json:"groups"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Split name into first and last name
	nameParts := strings.Fields(req.Name)
	firstName := ""
	lastName := ""
	if len(nameParts) > 0 {
		firstName = nameParts[0]
		if len(nameParts) > 1 {
			lastName = strings.Join(nameParts[1:], " ")
		}
	}

	user := models.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      models.UserRole(req.Role),
	}

	if len(req.Groups) > 0 {
		user.UserGroupID = &req.Groups[0] // TODO: Support multiple groups
	}

	if err := h.DB.Create(&user).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.WriteSuccess(w, user, "User created successfully")
}

func (h *AdminEnhancedHandler) UpdateAdminUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	var req struct {
		Name       string   `json:"name"`
		Email      string   `json:"email"`
		Role       string   `json:"role"`
		Department string   `json:"department"`
		Groups     []uint   `json:"groups"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Split name into first and last name
	nameParts := strings.Fields(req.Name)
	if len(nameParts) > 0 {
		user.FirstName = nameParts[0]
		if len(nameParts) > 1 {
			user.LastName = strings.Join(nameParts[1:], " ")
		}
	}
	user.Email = req.Email
	user.Role = models.UserRole(req.Role)

	if len(req.Groups) > 0 {
		user.UserGroupID = &req.Groups[0]
	}

	if err := h.DB.Save(&user).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	utils.WriteSuccess(w, user, "User updated successfully")
}

func (h *AdminEnhancedHandler) DeleteAdminUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.DB.Delete(&models.User{}, userID).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	utils.WriteSuccess(w, nil, "User deleted successfully")
}

// Enhanced Claim Types
type EnhancedClaimTypeRequest struct {
	Name                  string  `json:"name"`
	Code                  string  `json:"code"`
	Description           string  `json:"description"`
	Category              string  `json:"category"`
	MaxAmount             float64 `json:"maxAmount"`
	Icon                  string  `json:"icon"`
	Color                 string  `json:"color"`
	RequiresReceipt       bool    `json:"requiresReceipt"`
	RequiresApproval      bool    `json:"requiresApproval"`
	RequiresJustification bool    `json:"requiresJustification"`
	ApprovalLevels        int     `json:"approvalLevels"`
	ValidityPeriod        int     `json:"validityPeriod"`
	Active                bool    `json:"active"`
}

func (h *AdminEnhancedHandler) GetEnhancedClaimTypes(w http.ResponseWriter, r *http.Request) {
	var claimTypes []models.ClaimType
	if err := h.DB.Find(&claimTypes).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve claim types")
		return
	}

	// Convert to enhanced response
	type EnhancedClaimType struct {
		models.ClaimType
		Code                  string  `json:"code"`
		Category              string  `json:"category"`
		MaxAmount             float64 `json:"maxAmount"`
		Icon                  string  `json:"icon"`
		Color                 string  `json:"color"`
		RequiresReceipt       bool    `json:"requiresReceipt"`
		RequiresApproval      bool    `json:"requiresApproval"`
		RequiresJustification bool    `json:"requiresJustification"`
		ApprovalLevels        int     `json:"approvalLevels"`
		ValidityPeriod        int     `json:"validityPeriod"`
		Active                bool    `json:"active"`
	}

	var enhanced []EnhancedClaimType
	for _, ct := range claimTypes {
		enhanced = append(enhanced, EnhancedClaimType{
			ClaimType:             ct,
			Code:                  ct.Name, // TODO: Add code field to model
			Category:              "other",  // TODO: Add category field
			MaxAmount:             10000,    // TODO: Add max amount field
			Icon:                  "pi pi-tag",
			Color:                 "#3b82f6",
			RequiresReceipt:       true,
			RequiresApproval:      true,
			RequiresJustification: false,
			ApprovalLevels:        1,
			ValidityPeriod:        30,
			Active:                true,
		})
	}

	utils.WriteSuccess(w, enhanced)
}

func (h *AdminEnhancedHandler) CreateEnhancedClaimType(w http.ResponseWriter, r *http.Request) {
	var req EnhancedClaimTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	claimType := models.ClaimType{
		Name:        req.Name,
		Description: req.Description,
		// TODO: Add additional fields to ClaimType model
	}

	if err := h.DB.Create(&claimType).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create claim type")
		return
	}

	utils.WriteSuccess(w, claimType, "Claim type created successfully")
}

func (h *AdminEnhancedHandler) UpdateEnhancedClaimType(w http.ResponseWriter, r *http.Request) {
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

	var req EnhancedClaimTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	claimType.Name = req.Name
	claimType.Description = req.Description
	// TODO: Update additional fields

	if err := h.DB.Save(&claimType).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update claim type")
		return
	}

	utils.WriteSuccess(w, claimType, "Claim type updated successfully")
}

// Enhanced Groups
type EnhancedGroupRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Department  string   `json:"department"`
	Permissions []string `json:"permissions"`
	Members     []uint   `json:"members"`
}

func (h *AdminEnhancedHandler) GetEnhancedGroups(w http.ResponseWriter, r *http.Request) {
	var groups []models.UserGroup
	if err := h.DB.Find(&groups).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve groups")
		return
	}

	// Enhanced group response
	type EnhancedGroup struct {
		models.UserGroup
		Department  string        `json:"department"`
		Permissions []string      `json:"permissions"`
		Members     []models.User `json:"members"`
	}

	var enhanced []EnhancedGroup
	for _, group := range groups {
		// Get members for this group
		var members []models.User
		h.DB.Where("user_group_id = ?", group.ID).Find(&members)
		
		enhanced = append(enhanced, EnhancedGroup{
			UserGroup:   group,
			Department:  "all", // TODO: Add department field
			Permissions: []string{"claims.view", "claims.create"}, // TODO: Implement permissions
			Members:     members,
		})
	}

	utils.WriteSuccess(w, enhanced)
}

func (h *AdminEnhancedHandler) CreateEnhancedGroup(w http.ResponseWriter, r *http.Request) {
	var req EnhancedGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	group := models.UserGroup{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.DB.Create(&group).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create group")
		return
	}

	// Add members
	if len(req.Members) > 0 {
		var users []models.User
		h.DB.Find(&users, req.Members)
		for i := range users {
			users[i].UserGroupID = &group.ID
		}
		h.DB.Save(&users)
	}

	utils.WriteSuccess(w, group, "Group created successfully")
}

func (h *AdminEnhancedHandler) UpdateEnhancedGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	var group models.UserGroup
	if err := h.DB.First(&group, groupID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Group not found")
		return
	}

	var req EnhancedGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	group.Name = req.Name
	group.Description = req.Description

	if err := h.DB.Save(&group).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update group")
		return
	}

	// Update members
	if len(req.Members) > 0 {
		// Remove existing members
		h.DB.Model(&models.User{}).Where("user_group_id = ?", groupID).Update("user_group_id", nil)
		
		// Add new members
		var users []models.User
		h.DB.Find(&users, req.Members)
		for i := range users {
			users[i].UserGroupID = &group.ID
		}
		h.DB.Save(&users)
	}

	utils.WriteSuccess(w, group, "Group updated successfully")
}

func (h *AdminEnhancedHandler) DeleteEnhancedGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Remove users from group first
	h.DB.Model(&models.User{}).Where("user_group_id = ?", groupID).Update("user_group_id", nil)

	if err := h.DB.Delete(&models.UserGroup{}, groupID).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete group")
		return
	}

	utils.WriteSuccess(w, nil, "Group deleted successfully")
}

// Enhanced Approval Levels
type EnhancedApprovalLevelRequest struct {
	Name                 string      `json:"name"`
	Description          string      `json:"description"`
	MinAmount            float64     `json:"minAmount"`
	MaxAmount            *float64    `json:"maxAmount"`
	ClaimTypes           []string    `json:"claimTypes"`
	Approvers            []Approver  `json:"approvers"`
	RequiresAllApprovers bool        `json:"requiresAllApprovers"`
	AutoApprove          bool        `json:"autoApprove"`
	NotifyApprovers      bool        `json:"notifyApprovers"`
	EscalationDays       *int        `json:"escalationDays"`
	ReminderDays         *int        `json:"reminderDays"`
	// Status permissions
	CanDraft                bool        `json:"canDraft"`
	CanSubmit               bool        `json:"canSubmit"`
	CanApprove              bool        `json:"canApprove"`
	CanReject               bool        `json:"canReject"`
	CanSetPaymentInProgress bool        `json:"canSetPaymentInProgress"`
	CanSetPaid              bool        `json:"canSetPaid"`
}

type Approver struct {
	Type  string `json:"type"`  // "user", "group", "role"
	ID    uint   `json:"id"`
	Value string `json:"value"` // For role type
	Name  string `json:"name"`  // Display name
}

func (h *AdminEnhancedHandler) GetEnhancedApprovalLevels(w http.ResponseWriter, r *http.Request) {
	var levels []models.ApprovalLevel
	query := h.DB.Preload("UserGroup").Preload("Approver").Order("user_group_id, level")
	
	// Filter by user group if specified
	groupID := r.URL.Query().Get("groupId")
	if groupID != "" {
		query = query.Where("user_group_id = ?", groupID)
	}
	
	if err := query.Find(&levels).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve approval levels")
		return
	}

	// Enhanced response
	type EnhancedApprovalLevel struct {
		ID                   uint       `json:"id"`
		Level                int        `json:"level"`
		UserGroupID          uint       `json:"userGroupId"`
		UserGroup            *models.UserGroup `json:"userGroup,omitempty"`
		Name                 string     `json:"name"`
		Description          string     `json:"description"`
		MinAmount            float64    `json:"minAmount"`
		MaxAmount            *float64   `json:"maxAmount"`
		ClaimTypes           []string   `json:"claimTypes"`
		Approvers            []Approver `json:"approvers"`
		RequiresAllApprovers bool       `json:"requiresAllApprovers"`
		AutoApprove          bool       `json:"autoApprove"`
		NotifyApprovers      bool       `json:"notifyApprovers"`
		EscalationDays       *int       `json:"escalationDays"`
		ReminderDays         *int       `json:"reminderDays"`
		// Status permissions
		CanDraft                bool       `json:"canDraft"`
		CanSubmit               bool       `json:"canSubmit"`
		CanApprove              bool       `json:"canApprove"`
		CanReject               bool       `json:"canReject"`
		CanSetPaymentInProgress bool       `json:"canSetPaymentInProgress"`
		CanSetPaid              bool       `json:"canSetPaid"`
	}

	var enhanced []EnhancedApprovalLevel
	for _, level := range levels {
		approvers := []Approver{}
		if level.ApproverID > 0 && level.Approver.ID > 0 {
			approvers = append(approvers, Approver{
				Type: "user",
				ID:   level.Approver.ID,
				Name: level.Approver.FirstName + " " + level.Approver.LastName,
			})
		}

		levelName := "Level " + strconv.Itoa(level.Level)
		if level.UserGroup.Name != "" {
			levelName += " - " + level.UserGroup.Name
		}
		
		enhanced = append(enhanced, EnhancedApprovalLevel{
			ID:                   level.ID,
			Level:                level.Level,
			UserGroupID:          level.UserGroupID,
			UserGroup:            &level.UserGroup,
			Name:                 levelName,
			Description:          "Approval level " + strconv.Itoa(level.Level) + " for " + level.UserGroup.Name,
			MinAmount:            0,
			MaxAmount:            nil,
			ClaimTypes:           []string{}, // All types
			Approvers:            approvers,
			RequiresAllApprovers: false,
			AutoApprove:          false,
			NotifyApprovers:      true,
			EscalationDays:       nil,
			ReminderDays:         nil,
			// Status permissions from database
			CanDraft:                level.CanDraft,
			CanSubmit:               level.CanSubmit,
			CanApprove:              level.CanApprove,
			CanReject:               level.CanReject,
			CanSetPaymentInProgress: level.CanSetPaymentInProgress,
			CanSetPaid:              level.CanSetPaid,
		})
	}

	utils.WriteSuccess(w, enhanced)
}

func (h *AdminEnhancedHandler) CreateEnhancedApprovalLevel(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserGroupID uint `json:"userGroupId"`
		ApproverID  uint `json:"approverId"`
		// Status permissions
		CanDraft                bool `json:"canDraft"`
		CanSubmit               bool `json:"canSubmit"`
		CanApprove              bool `json:"canApprove"`
		CanReject               bool `json:"canReject"`
		CanSetPaymentInProgress bool `json:"canSetPaymentInProgress"`
		CanSetPaid              bool `json:"canSetPaid"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate user group exists
	var userGroup models.UserGroup
	if err := h.DB.First(&userGroup, req.UserGroupID).Error; err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user group")
		return
	}

	// Validate approver exists
	var approver models.User
	if err := h.DB.First(&approver, req.ApproverID).Error; err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid approver")
		return
	}

	// Get the next level number for this user group
	var maxLevel int
	h.DB.Model(&models.ApprovalLevel{}).
		Where("user_group_id = ?", req.UserGroupID).
		Select("COALESCE(MAX(level), 0)").
		Scan(&maxLevel)

	level := models.ApprovalLevel{
		Level:                   maxLevel + 1,
		UserGroupID:             req.UserGroupID,
		ApproverID:              req.ApproverID,
		CanDraft:                req.CanDraft,
		CanSubmit:               req.CanSubmit,
		CanApprove:              req.CanApprove,
		CanReject:               req.CanReject,
		CanSetPaymentInProgress: req.CanSetPaymentInProgress,
		CanSetPaid:              req.CanSetPaid,
	}

	if err := h.DB.Create(&level).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create approval level")
		return
	}

	// Load associations for response
	h.DB.Preload("UserGroup").Preload("Approver").First(&level, level.ID)

	utils.WriteSuccess(w, level, "Approval level created successfully")
}

func (h *AdminEnhancedHandler) UpdateEnhancedApprovalLevel(w http.ResponseWriter, r *http.Request) {
	levelID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid approval level ID")
		return
	}

	var level models.ApprovalLevel
	if err := h.DB.First(&level, levelID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Approval level not found")
		return
	}

	var req EnhancedApprovalLevelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update status permissions
	level.CanDraft = req.CanDraft
	level.CanSubmit = req.CanSubmit
	level.CanApprove = req.CanApprove
	level.CanReject = req.CanReject
	level.CanSetPaymentInProgress = req.CanSetPaymentInProgress
	level.CanSetPaid = req.CanSetPaid

	// Update approver if provided
	if len(req.Approvers) > 0 && req.Approvers[0].Type == "user" {
		level.ApproverID = req.Approvers[0].ID
	}

	if err := h.DB.Save(&level).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update approval level")
		return
	}

	utils.WriteSuccess(w, level, "Approval level updated successfully")
}

func (h *AdminEnhancedHandler) DeleteEnhancedApprovalLevel(w http.ResponseWriter, r *http.Request) {
	levelID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid approval level ID")
		return
	}

	// Get the level to be deleted
	var levelToDelete models.ApprovalLevel
	if err := h.DB.First(&levelToDelete, levelID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Approval level not found")
		return
	}

	// Delete the level
	if err := h.DB.Delete(&models.ApprovalLevel{}, levelID).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete approval level")
		return
	}

	// Reorder remaining levels for this user group
	var remainingLevels []models.ApprovalLevel
	h.DB.Where("user_group_id = ? AND level > ?", levelToDelete.UserGroupID, levelToDelete.Level).
		Order("level").Find(&remainingLevels)
	
	// Update levels to be continuous
	for i, level := range remainingLevels {
		newLevel := levelToDelete.Level + i
		h.DB.Model(&level).Update("level", newLevel)
	}

	utils.WriteSuccess(w, nil, "Approval level deleted successfully")
}

func (h *AdminEnhancedHandler) UpdateApprovalLevelOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserGroupID uint `json:"userGroupId"`
		Orders []struct {
			ID    uint `json:"id"`
			Level int  `json:"level"`
		} `json:"orders"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Begin transaction
	tx := h.DB.Begin()

	// Update each level
	for _, order := range req.Orders {
		if err := tx.Model(&models.ApprovalLevel{}).
			Where("id = ? AND user_group_id = ?", order.ID, req.UserGroupID).
			Update("level", order.Level).Error; err != nil {
			tx.Rollback()
			utils.WriteError(w, http.StatusInternalServerError, "Failed to update approval level order")
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update approval level order")
		return
	}

	utils.WriteSuccess(w, nil, "Approval level order updated successfully")
}

// GetApprovalLevelsByGroup returns approval levels grouped by user group
func (h *AdminEnhancedHandler) GetApprovalLevelsByGroup(w http.ResponseWriter, r *http.Request) {
	var groups []models.UserGroup
	if err := h.DB.Find(&groups).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve user groups")
		return
	}

	type GroupApprovalLevels struct {
		GroupID   uint   `json:"groupId"`
		GroupName string `json:"groupName"`
		Levels    []struct {
			ID         uint   `json:"id"`
			Level      int    `json:"level"`
			ApproverID uint   `json:"approverId"`
			Approver   struct {
				ID        uint   `json:"id"`
				Name      string `json:"name"`
				Email     string `json:"email"`
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
			} `json:"approver"`
			Permissions struct {
				CanDraft                bool `json:"canDraft"`
				CanSubmit               bool `json:"canSubmit"`
				CanApprove              bool `json:"canApprove"`
				CanReject               bool `json:"canReject"`
				CanSetPaymentInProgress bool `json:"canSetPaymentInProgress"`
				CanSetPaid              bool `json:"canSetPaid"`
			} `json:"permissions"`
		} `json:"levels"`
	}

	var result []GroupApprovalLevels

	for _, group := range groups {
		var levels []models.ApprovalLevel
		h.DB.Preload("Approver").Where("user_group_id = ?", group.ID).Order("level").Find(&levels)

		groupLevels := GroupApprovalLevels{
			GroupID:   group.ID,
			GroupName: group.Name,
			Levels:    make([]struct {
				ID         uint   `json:"id"`
				Level      int    `json:"level"`
				ApproverID uint   `json:"approverId"`
				Approver   struct {
					ID        uint   `json:"id"`
					Name      string `json:"name"`
					Email     string `json:"email"`
					FirstName string `json:"firstName"`
					LastName  string `json:"lastName"`
				} `json:"approver"`
				Permissions struct {
					CanDraft                bool `json:"canDraft"`
					CanSubmit               bool `json:"canSubmit"`
					CanApprove              bool `json:"canApprove"`
					CanReject               bool `json:"canReject"`
					CanSetPaymentInProgress bool `json:"canSetPaymentInProgress"`
					CanSetPaid              bool `json:"canSetPaid"`
				} `json:"permissions"`
			}, 0),
		}

		for _, level := range levels {
			levelData := struct {
				ID         uint   `json:"id"`
				Level      int    `json:"level"`
				ApproverID uint   `json:"approverId"`
				Approver   struct {
					ID        uint   `json:"id"`
					Name      string `json:"name"`
					Email     string `json:"email"`
					FirstName string `json:"firstName"`
					LastName  string `json:"lastName"`
				} `json:"approver"`
				Permissions struct {
					CanDraft                bool `json:"canDraft"`
					CanSubmit               bool `json:"canSubmit"`
					CanApprove              bool `json:"canApprove"`
					CanReject               bool `json:"canReject"`
					CanSetPaymentInProgress bool `json:"canSetPaymentInProgress"`
					CanSetPaid              bool `json:"canSetPaid"`
				} `json:"permissions"`
			}{
				ID:         level.ID,
				Level:      level.Level,
				ApproverID: level.ApproverID,
			}

			if level.Approver.ID > 0 {
				levelData.Approver.ID = level.Approver.ID
				levelData.Approver.Name = level.Approver.FirstName + " " + level.Approver.LastName
				levelData.Approver.Email = level.Approver.Email
				levelData.Approver.FirstName = level.Approver.FirstName
				levelData.Approver.LastName = level.Approver.LastName
			}

			levelData.Permissions.CanDraft = level.CanDraft
			levelData.Permissions.CanSubmit = level.CanSubmit
			levelData.Permissions.CanApprove = level.CanApprove
			levelData.Permissions.CanReject = level.CanReject
			levelData.Permissions.CanSetPaymentInProgress = level.CanSetPaymentInProgress
			levelData.Permissions.CanSetPaid = level.CanSetPaid

			groupLevels.Levels = append(groupLevels.Levels, levelData)
		}

		result = append(result, groupLevels)
	}

	utils.WriteSuccess(w, result)
}