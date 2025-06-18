package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"hrcs/backend/middleware"
	"hrcs/backend/models"
	"hrcs/backend/services"
	"hrcs/backend/utils"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type ClaimHandler struct {
	DB             *gorm.DB
	balanceService *services.BalanceService
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
	return &ClaimHandler{
		DB:             db,
		balanceService: services.NewBalanceService(db),
	}
}

func (h *ClaimHandler) GetClaimTypes(w http.ResponseWriter, r *http.Request) {
	var claimTypes []models.ClaimType
	if err := h.DB.Find(&claimTypes).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve claim types")
		return
	}

	utils.WriteSuccess(w, claimTypes)
}

func (h *ClaimHandler) CreateClaim(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	
	var req CreateClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// NEW: Balance validation
	canClaim, message, err := h.balanceService.CanUserClaim(user.ID, req.ClaimTypeID, req.Amount)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Balance check failed")
		return
	}
	if !canClaim {
		utils.WriteError(w, http.StatusBadRequest, message)
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

	// Store the old status to check if we're changing to paid
	oldStatus := claim.Status
	claim.Status = req.Status
	
	if err := h.DB.Save(&claim).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update claim status")
		return
	}

	// NEW: Deduct from balance when marking as paid
	if req.Status == models.StatusPaid && oldStatus != models.StatusPaid {
		err := h.balanceService.DeductFromBalance(claim.UserID, claim.ClaimTypeID, claim.Amount)
		if err != nil {
			// Rollback the status change
			claim.Status = oldStatus
			h.DB.Save(&claim)
			utils.WriteError(w, http.StatusInternalServerError, "Balance deduction failed")
			return
		}
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

// NEW: Enhanced workflow endpoints for Phase 3

// GetClaimWorkflow returns detailed workflow information for a claim
func (h *ClaimHandler) GetClaimWorkflow(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	var claim models.Claim
	query := h.DB.Preload("User").Preload("ClaimType").Preload("Approvals.Approver").Preload("Approvals.ApprovalLevel.UserGroup")

	// Check permissions
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

	// Get all approval levels for workflow structure
	var approvalLevels []models.ApprovalLevel
	h.DB.Preload("UserGroup").Preload("Approver").Order("level ASC").Find(&approvalLevels)

	// Build workflow response
	workflow := buildWorkflowResponse(&claim, approvalLevels)

	utils.WriteSuccess(w, workflow)
}

// UpdateClaimWorkflowStep updates a specific workflow step
func (h *ClaimHandler) UpdateClaimWorkflowStep(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	claimID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	stepID, err := strconv.Atoi(chi.URLParam(r, "stepId"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid step ID")
		return
	}

	var req struct {
		Action   string `json:"action"`   // approve, reject, request_changes
		Comments string `json:"comments"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get the claim
	var claim models.Claim
	if err := h.DB.First(&claim, claimID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Claim not found")
		return
	}

	// Check if user can approve this step
	var approvalLevel models.ApprovalLevel
	if err := h.DB.Where("id = ? AND approver_id = ?", stepID, user.ID).First(&approvalLevel).Error; err != nil {
		utils.WriteError(w, http.StatusForbidden, "You are not authorized to approve this step")
		return
	}

	// Check permissions based on action
	canPerformAction := false
	var newStatus models.ClaimStatus

	switch req.Action {
	case "approve":
		canPerformAction = approvalLevel.CanApprove
		newStatus = models.StatusApproved
	case "reject":
		canPerformAction = approvalLevel.CanReject
		newStatus = models.StatusRejected
	case "request_changes":
		canPerformAction = approvalLevel.CanReject
		newStatus = models.StatusDraft
	default:
		utils.WriteError(w, http.StatusBadRequest, "Invalid action")
		return
	}

	if !canPerformAction {
		utils.WriteError(w, http.StatusForbidden, "You don't have permission to perform this action")
		return
	}

	// Update claim status
	oldStatus := claim.Status
	claim.Status = newStatus
	
	if err := h.DB.Save(&claim).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update claim status")
		return
	}

	// Handle balance deduction for approved claims
	if newStatus == models.StatusApproved && oldStatus != models.StatusApproved {
		err := h.balanceService.DeductFromBalance(claim.UserID, claim.ClaimTypeID, claim.Amount)
		if err != nil {
			// Rollback the status change
			claim.Status = oldStatus
			h.DB.Save(&claim)
			utils.WriteError(w, http.StatusInternalServerError, "Balance deduction failed")
			return
		}
	}

	// Create approval record
	approval := models.ClaimApproval{
		ClaimID:         claim.ID,
		ApprovalLevelID: uint(stepID),
		ApproverID:      user.ID,
		Status:          newStatus,
		Comments:        req.Comments,
	}

	if err := h.DB.Create(&approval).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create approval record")
		return
	}

	// Return updated workflow
	h.GetClaimWorkflow(w, r)
}

// Helper function to build workflow response
func buildWorkflowResponse(claim *models.Claim, approvalLevels []models.ApprovalLevel) map[string]interface{} {
	// Build workflow steps
	var workflowSteps []map[string]interface{}

	// Add creation step
	creationStep := map[string]interface{}{
		"id":          "created",
		"type":        "system",
		"title":       "Claim Created",
		"description": "Claim created by " + claim.User.FirstName + " " + claim.User.LastName,
		"status":      "completed",
		"completed_at": claim.CreatedAt,
		"level":       0,
		"user":        claim.User,
	}
	workflowSteps = append(workflowSteps, creationStep)

	// Add submission step if submitted
	if claim.Status != models.StatusDraft {
		submissionStep := map[string]interface{}{
			"id":          "submitted",
			"type":        "system", 
			"title":       "Submitted for Approval",
			"description": "Claim submitted to approval workflow",
			"status":      "completed",
			"completed_at": claim.UpdatedAt, // Approximation
			"level":       0,
		}
		workflowSteps = append(workflowSteps, submissionStep)
	}

	// Add approval level steps
	for _, level := range approvalLevels {
		// Find corresponding approval record
		var approval *models.ClaimApproval
		for _, a := range claim.Approvals {
			if a.ApprovalLevelID == level.ID {
				approval = &a
				break
			}
		}

		step := map[string]interface{}{
			"id":            level.ID,
			"type":          "approval",
			"title":         level.UserGroup.Name + " Approval",
			"description":   "Level " + strconv.Itoa(level.Level) + " approval required",
			"level":         level.Level,
			"approver":      level.Approver,
			"user_group":    level.UserGroup,
			"permissions":   getApprovalPermissions(level),
		}

		if approval != nil {
			step["status"] = string(approval.Status)
			step["completed_at"] = approval.CreatedAt
			step["comments"] = approval.Comments
			step["approved_by"] = approval.Approver
		} else {
			// Determine step status based on claim status and level
			if claim.Status == models.StatusDraft {
				step["status"] = "not_started"
			} else if claim.Status == models.StatusSubmitted {
				step["status"] = "pending"
			} else {
				step["status"] = "not_reached"
			}
		}

		workflowSteps = append(workflowSteps, step)
	}

	// Calculate progress
	completedSteps := 0
	totalSteps := len(workflowSteps)
	
	for _, step := range workflowSteps {
		if status, ok := step["status"].(string); ok && (status == "completed" || status == "approved") {
			completedSteps++
		}
	}

	progress := float64(completedSteps) / float64(totalSteps) * 100

	// Determine current step
	var currentStep map[string]interface{}
	for _, step := range workflowSteps {
		if status, ok := step["status"].(string); ok && status == "pending" {
			currentStep = step
			break
		}
	}

	// Determine next possible actions
	nextActions := getNextActions(claim, approvalLevels)

	return map[string]interface{}{
		"claim":         claim,
		"workflow_steps": workflowSteps,
		"progress":      progress,
		"completed_steps": completedSteps,
		"total_steps":   totalSteps,
		"current_step":  currentStep,
		"next_actions":  nextActions,
		"status":        claim.Status,
	}
}

// Helper function to get approval permissions
func getApprovalPermissions(level models.ApprovalLevel) map[string]bool {
	return map[string]bool{
		"can_approve":                level.CanApprove,
		"can_reject":                level.CanReject,
		"can_set_payment_in_progress": level.CanSetPaymentInProgress,
		"can_set_paid":              level.CanSetPaid,
		"can_draft":                 level.CanDraft,
		"can_submit":                level.CanSubmit,
	}
}

// Helper function to get next possible actions
func getNextActions(claim *models.Claim, approvalLevels []models.ApprovalLevel) []map[string]interface{} {
	var actions []map[string]interface{}

	switch claim.Status {
	case models.StatusDraft:
		actions = append(actions, map[string]interface{}{
			"action":      "submit",
			"label":       "Submit for Approval",
			"description": "Submit claim to approval workflow",
			"available":   true,
		})
		actions = append(actions, map[string]interface{}{
			"action":      "edit",
			"label":       "Edit Claim",
			"description": "Modify claim details",
			"available":   true,
		})
		actions = append(actions, map[string]interface{}{
			"action":      "cancel",
			"label":       "Cancel Claim",
			"description": "Delete this claim",
			"available":   true,
		})

	case models.StatusSubmitted:
		// Find next approval level that hasn't been completed
		for _, level := range approvalLevels {
			hasApproval := false
			for _, approval := range claim.Approvals {
				if approval.ApprovalLevelID == level.ID {
					hasApproval = true
					break
				}
			}
			
			if !hasApproval {
				if level.CanApprove {
					actions = append(actions, map[string]interface{}{
						"action":      "approve",
						"label":       "Approve",
						"description": "Approve this claim",
						"level_id":    level.ID,
						"available":   true,
					})
				}
				if level.CanReject {
					actions = append(actions, map[string]interface{}{
						"action":      "reject",
						"label":       "Reject",
						"description": "Reject this claim",
						"level_id":    level.ID,
						"available":   true,
					})
					actions = append(actions, map[string]interface{}{
						"action":      "request_changes",
						"label":       "Request Changes",
						"description": "Request changes to this claim",
						"level_id":    level.ID,
						"available":   true,
					})
				}
				break // Only show actions for the next level
			}
		}

	case models.StatusApproved:
		// Look for levels that can set payment status
		for _, level := range approvalLevels {
			if level.CanSetPaymentInProgress {
				actions = append(actions, map[string]interface{}{
					"action":      "set_payment_in_progress",
					"label":       "Mark Payment in Progress",
					"description": "Mark claim as payment in progress",
					"level_id":    level.ID,
					"available":   true,
				})
			}
			if level.CanSetPaid {
				actions = append(actions, map[string]interface{}{
					"action":      "set_paid",
					"label":       "Mark as Paid",
					"description": "Mark claim as paid",
					"level_id":    level.ID,
					"available":   true,
				})
			}
		}
	}

	return actions
}