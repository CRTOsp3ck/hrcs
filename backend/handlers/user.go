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

type UserHandler struct {
	DB *gorm.DB
}

type UpdateUserRoleRequest struct {
	Role models.UserRole `json:"role"`
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	utils.WriteSuccess(w, user)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := h.DB.Preload("UserGroup").Find(&users).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	utils.WriteSuccess(w, users)
}

func (h *UserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req UpdateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	user.Role = req.Role
	if err := h.DB.Save(&user).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update user role")
		return
	}

	utils.WriteSuccess(w, user, "User role updated successfully")
}