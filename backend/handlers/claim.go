package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"hrcs/backend/middleware"
	"hrcs/backend/models"
	"hrcs/backend/utils"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type ClaimHandler struct {
	DB *gorm.DB
}

type CreateClaimRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	ClaimTypeID uint    `json:"claim_type_id"`
}

type UpdateClaimRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	ClaimTypeID uint    `json:"claim_type_id"`
}

type ApproveClaimRequest struct {
	Status   models.ClaimStatus `json:"status"`
	Comments string             `json:"comments"`
}

func NewClaimHandler(db *gorm.DB) *ClaimHandler {
	return &ClaimHandler{DB: db}
}

func (h *ClaimHandler) CreateClaim(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	
	var req CreateClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	claim := models.Claim{
		Title:       req.Title,
		Description: req.Description,
		Amount:      req.Amount,
		UserID:      user.ID,
		ClaimTypeID: req.ClaimTypeID,
		Status:      models.StatusDraft,
	}

	if err := h.DB.Create(&claim).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create claim")
		return
	}

	if err := h.DB.Preload("User").Preload("ClaimType").First(&claim, claim.ID).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve claim")
		return
	}

	utils.WriteSuccess(w, claim, "Claim created successfully")
}

func (h *ClaimHandler) GetClaims(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	
	var claims []models.Claim
	query := h.DB.Preload("User").Preload("ClaimType").Preload("Approvals")

	if user.Role == models.RoleAdmin {
		query = query.Find(&claims)
	} else {
		query = query.Where("user_id = ?", user.ID).Find(&claims)
	}

	if err := query.Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve claims")
		return
	}

	utils.WriteSuccess(w, claims)
}

func (h *ClaimHandler) GetClaim(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	var claim models.Claim
	query := h.DB.Preload("User").Preload("ClaimType").Preload("Approvals.Approver").Preload("Approvals.ApprovalLevel")

	if user.Role == models.RoleAdmin {
		query = query.First(&claim, claimID)
	} else {
		query = query.Where("user_id = ?", user.ID).First(&claim, claimID)
	}

	if err := query.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.WriteError(w, http.StatusNotFound, "Claim not found")
		} else {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve claim")
		}
		return
	}

	utils.WriteSuccess(w, claim)
}

func (h *ClaimHandler) UpdateClaim(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	var claim models.Claim
	if err := h.DB.Where("user_id = ? AND id = ?", user.ID, claimID).First(&claim).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Claim not found")
		return
	}

	if claim.Status != models.StatusDraft {
		utils.WriteError(w, http.StatusBadRequest, "Can only update draft claims")
		return
	}

	var req UpdateClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	claim.Title = req.Title
	claim.Description = req.Description
	claim.Amount = req.Amount
	claim.ClaimTypeID = req.ClaimTypeID

	if err := h.DB.Save(&claim).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update claim")
		return
	}

	utils.WriteSuccess(w, claim, "Claim updated successfully")
}

func (h *ClaimHandler) SubmitClaim(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	var claim models.Claim
	if err := h.DB.Where("user_id = ? AND id = ?", user.ID, claimID).First(&claim).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Claim not found")
		return
	}

	if claim.Status != models.StatusDraft {
		utils.WriteError(w, http.StatusBadRequest, "Can only submit draft claims")
		return
	}

	claim.Status = models.StatusSubmitted
	if err := h.DB.Save(&claim).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to submit claim")
		return
	}

	utils.WriteSuccess(w, claim, "Claim submitted successfully")
}

func (h *ClaimHandler) CancelClaim(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	var claim models.Claim
	if err := h.DB.Where("user_id = ? AND id = ?", user.ID, claimID).First(&claim).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Claim not found")
		return
	}

	if claim.Status == models.StatusApproved || claim.Status == models.StatusPaid || claim.Status == models.StatusPaymentInProgress {
		utils.WriteError(w, http.StatusBadRequest, "Cannot cancel approved or processed claims")
		return
	}

	if err := h.DB.Delete(&claim).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to cancel claim")
		return
	}

	utils.WriteSuccess(w, nil, "Claim cancelled successfully")
}

func (h *ClaimHandler) ApproveClaim(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	var claim models.Claim
	if err := h.DB.Preload("User").First(&claim, claimID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Claim not found")
		return
	}

	if claim.UserID == user.ID {
		utils.WriteError(w, http.StatusForbidden, "Cannot approve your own claim")
		return
	}

	var req ApproveClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	claim.Status = req.Status
	if err := h.DB.Save(&claim).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update claim status")
		return
	}

	approval := models.ClaimApproval{
		ClaimID:    claim.ID,
		ApproverID: user.ID,
		Status:     req.Status,
		Comments:   req.Comments,
	}

	if err := h.DB.Create(&approval).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create approval record")
		return
	}

	utils.WriteSuccess(w, claim, "Claim status updated successfully")
}