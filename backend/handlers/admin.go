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
	Level       int  `json:"level"`
	UserGroupID uint `json:"user_group_id"`
	ApproverID  uint `json:"approver_id"`
	CanApprove  bool `json:"can_approve"`
	CanReject   bool `json:"can_reject"`
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
		Level:       req.Level,
		UserGroupID: req.UserGroupID,
		ApproverID:  req.ApproverID,
		CanApprove:  req.CanApprove,
		CanReject:   req.CanReject,
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