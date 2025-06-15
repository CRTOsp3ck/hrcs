package handlers

import (
	"encoding/json"
	"hrcs/backend/services"
	"hrcs/backend/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type BalanceHandler struct {
	balanceService *services.BalanceService
}

func NewBalanceHandler(db *gorm.DB) *BalanceHandler {
	return &BalanceHandler{
		balanceService: services.NewBalanceService(db),
	}
}

// GetUserBalances gets all balance records for the authenticated user
func (h *BalanceHandler) GetUserBalances(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	balances, err := h.balanceService.GetAllUserBalances(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get user balances")
		return
	}

	utils.WriteSuccess(w, map[string]interface{}{
		"balances": balances,
	})
}

// GetUserBalance gets balance for a specific claim type for the authenticated user
func (h *BalanceHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := chi.URLParam(r, "id")
	claimTypeID, err := strconv.ParseUint(vars, 10, 32)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid claim type ID")
		return
	}

	balance, err := h.balanceService.GetUserBalance(userID, uint(claimTypeID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get user balance")
		return
	}

	utils.WriteSuccess(w, map[string]interface{}{
		"balance": balance,
	})
}

// CheckClaimAmount validates if user can claim a specific amount
func (h *BalanceHandler) CheckClaimAmount(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req struct {
		ClaimTypeID uint    `json:"claim_type_id"`
		Amount      float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	canClaim, message, err := h.balanceService.CanUserClaim(userID, req.ClaimTypeID, req.Amount)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to check claim amount")
		return
	}

	utils.WriteSuccess(w, map[string]interface{}{
		"can_claim": canClaim,
		"message":   message,
	})
}

// AdminUpdateBalance allows admin to manually adjust a user's balance
func (h *BalanceHandler) AdminUpdateBalance(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	userRole, ok := r.Context().Value("user_role").(string)
	if !ok || userRole != "admin" {
		utils.WriteError(w, http.StatusForbidden, "Admin access required")
		return
	}

	var req struct {
		UserID      uint    `json:"user_id"`
		ClaimTypeID uint    `json:"claim_type_id"`
		NewLimit    float64 `json:"new_limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.balanceService.AdminAdjustBalance(req.UserID, req.ClaimTypeID, req.NewLimit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to adjust balance")
		return
	}

	utils.WriteSuccess(w, map[string]interface{}{
		"message": "Balance adjusted successfully",
	})
}

// GetUserBalanceDetails gets detailed balance information for admin view
func (h *BalanceHandler) GetUserBalanceDetails(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	userRole, ok := r.Context().Value("user_role").(string)
	if !ok || userRole != "admin" {
		utils.WriteError(w, http.StatusForbidden, "Admin access required")
		return
	}

	vars := chi.URLParam(r, "id")
	userID, err := strconv.ParseUint(vars, 10, 32)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	balances, err := h.balanceService.GetAllUserBalances(uint(userID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get user balance details")
		return
	}

	utils.WriteSuccess(w, map[string]interface{}{
		"user_id":  userID,
		"balances": balances,
	})
}