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

type AuditHandler struct {
	DB *gorm.DB
}

func NewAuditHandler(db *gorm.DB) *AuditHandler {
	return &AuditHandler{DB: db}
}

// GetAuditLogs returns paginated audit logs with filtering
func (h *AuditHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 25
	}

	// Filter parameters
	action := r.URL.Query().Get("action")
	entityType := r.URL.Query().Get("entity_type")
	userID := r.URL.Query().Get("user_id")
	search := r.URL.Query().Get("search")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// Build query
	query := h.DB.Model(&models.AuditLog{}).Preload("User")

	// Apply filters
	if action != "" {
		query = query.Where("action = ?", action)
	}

	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}

	if userID != "" {
		if uid, err := strconv.Atoi(userID); err == nil {
			query = query.Where("user_id = ?", uid)
		}
	}

	if search != "" {
		// Search in action, entity type, and user names
		query = query.Joins("LEFT JOIN users ON audit_logs.user_id = users.id").
			Where("audit_logs.action ILIKE ? OR audit_logs.entity_type ILIKE ? OR users.first_name ILIKE ? OR users.last_name ILIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("audit_logs.created_at >= ?", start)
		}
	}

	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			// Add 24 hours to include the entire end date
			endTime := end.Add(24 * time.Hour)
			query = query.Where("audit_logs.created_at < ?", endTime)
		}
	}

	// Get total count
	var total int64
	countQuery := *query
	countQuery.Count(&total)

	// Get paginated results
	offset := (page - 1) * limit
	var auditLogs []models.AuditLog
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&auditLogs).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve audit logs")
		return
	}

	// Transform to response format
	var responses []models.AuditLogResponse
	for _, log := range auditLogs {
		response := models.AuditLogResponse{
			ID:         log.ID,
			UserID:     log.UserID,
			UserName:   log.User.FirstName + " " + log.User.LastName,
			UserEmail:  log.User.Email,
			Action:     log.Action,
			EntityType: log.EntityType,
			EntityID:   log.EntityID,
			OldValues:  log.OldValues,
			NewValues:  log.NewValues,
			IPAddress:  log.IPAddress,
			UserAgent:  log.UserAgent,
			CreatedAt:  log.CreatedAt,
			Description: log.GetDescription(),
			ActionType: log.GetActionType(),
			Severity:   log.GetSeverity(),
		}
		responses = append(responses, response)
	}

	// Prepare pagination info
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	
	result := map[string]interface{}{
		"data": responses,
		"pagination": map[string]interface{}{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
			"has_next":     page < totalPages,
			"has_prev":     page > 1,
		},
		"filters": map[string]interface{}{
			"action":      action,
			"entity_type": entityType,
			"user_id":     userID,
			"search":      search,
			"start_date":  startDate,
			"end_date":    endDate,
		},
	}

	utils.WriteSuccess(w, result)
}

// GetAuditLogDetails returns detailed information about a specific audit log entry
func (h *AuditHandler) GetAuditLogDetails(w http.ResponseWriter, r *http.Request) {
	auditLogID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid audit log ID")
		return
	}

	var auditLog models.AuditLog
	if err := h.DB.Preload("User").First(&auditLog, auditLogID).Error; err != nil {
		utils.WriteError(w, http.StatusNotFound, "Audit log not found")
		return
	}

	// Parse old and new values if they exist
	var oldValues, newValues interface{}
	if auditLog.OldValues != "" {
		json.Unmarshal([]byte(auditLog.OldValues), &oldValues)
	}
	if auditLog.NewValues != "" {
		json.Unmarshal([]byte(auditLog.NewValues), &newValues)
	}

	response := map[string]interface{}{
		"id":          auditLog.ID,
		"user":        auditLog.User,
		"action":      auditLog.Action,
		"entity_type": auditLog.EntityType,
		"entity_id":   auditLog.EntityID,
		"old_values":  oldValues,
		"new_values":  newValues,
		"ip_address":  auditLog.IPAddress,
		"user_agent":  auditLog.UserAgent,
		"created_at":  auditLog.CreatedAt,
		"description": auditLog.GetDescription(),
		"action_type": auditLog.GetActionType(),
		"severity":    auditLog.GetSeverity(),
	}

	utils.WriteSuccess(w, response)
}

// GetAuditLogStats returns statistics about audit log activity
func (h *AuditHandler) GetAuditLogStats(w http.ResponseWriter, r *http.Request) {
	// Get date range from query params, default to last 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	if start := r.URL.Query().Get("start_date"); start != "" {
		if parsed, err := time.Parse("2006-01-02", start); err == nil {
			startDate = parsed
		}
	}

	if end := r.URL.Query().Get("end_date"); end != "" {
		if parsed, err := time.Parse("2006-01-02", end); err == nil {
			endDate = parsed.Add(24 * time.Hour) // Include the entire end date
		}
	}

	// Total activity count
	var totalActivity int64
	h.DB.Model(&models.AuditLog{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Count(&totalActivity)

	// Activity by action
	type ActionCount struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}
	var actionCounts []ActionCount
	h.DB.Model(&models.AuditLog{}).
		Select("action, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("action").
		Order("count DESC").
		Scan(&actionCounts)

	// Activity by entity type
	type EntityCount struct {
		EntityType string `json:"entity_type"`
		Count      int64  `json:"count"`
	}
	var entityCounts []EntityCount
	h.DB.Model(&models.AuditLog{}).
		Select("entity_type, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("entity_type").
		Order("count DESC").
		Scan(&entityCounts)

	// Top active users
	type UserActivity struct {
		UserID    uint   `json:"user_id"`
		UserName  string `json:"user_name"`
		UserEmail string `json:"user_email"`
		Count     int64  `json:"count"`
	}
	var userActivity []UserActivity
	h.DB.Table("audit_logs").
		Select("audit_logs.user_id, users.first_name || ' ' || users.last_name as user_name, users.email as user_email, COUNT(*) as count").
		Joins("LEFT JOIN users ON audit_logs.user_id = users.id").
		Where("audit_logs.created_at BETWEEN ? AND ?", startDate, endDate).
		Group("audit_logs.user_id, users.first_name, users.last_name, users.email").
		Order("count DESC").
		Limit(10).
		Scan(&userActivity)

	// Daily activity for the period
	type DailyActivity struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	var dailyActivity []DailyActivity
	h.DB.Model(&models.AuditLog{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&dailyActivity)

	stats := map[string]interface{}{
		"period": map[string]interface{}{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
		"total_activity":  totalActivity,
		"by_action":       actionCounts,
		"by_entity_type":  entityCounts,
		"top_users":       userActivity,
		"daily_activity":  dailyActivity,
	}

	utils.WriteSuccess(w, stats)
}

// CreateAuditLog creates a new audit log entry (internal use)
func (h *AuditHandler) CreateAuditLog(req models.AuditLogRequest) error {
	// Convert old and new values to JSON strings
	var oldValuesJSON, newValuesJSON string
	
	if req.OldValues != nil {
		if bytes, err := json.Marshal(req.OldValues); err == nil {
			oldValuesJSON = string(bytes)
		}
	}
	
	if req.NewValues != nil {
		if bytes, err := json.Marshal(req.NewValues); err == nil {
			newValuesJSON = string(bytes)
		}
	}

	auditLog := models.AuditLog{
		UserID:     req.UserID,
		Action:     req.Action,
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		OldValues:  oldValuesJSON,
		NewValues:  newValuesJSON,
		IPAddress:  req.IPAddress,
		UserAgent:  req.UserAgent,
		CreatedAt:  time.Now(),
	}

	return h.DB.Create(&auditLog).Error
}

// Helper function to get user IP address from request
func GetUserIP(r *http.Request) string {
	// Check for X-Forwarded-For header first (in case of proxy)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check for X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// Helper function to create an audit log entry from HTTP request context
func (h *AuditHandler) LogActivity(userID uint, action, entityType string, entityID uint, oldValues, newValues interface{}, r *http.Request) error {
	req := models.AuditLogRequest{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		IPAddress:  GetUserIP(r),
		UserAgent:  r.UserAgent(),
	}

	if oldValues != nil {
		if oldMap, ok := oldValues.(map[string]interface{}); ok {
			req.OldValues = oldMap
		}
	}

	if newValues != nil {
		if newMap, ok := newValues.(map[string]interface{}); ok {
			req.NewValues = newMap
		}
	}

	return h.CreateAuditLog(req)
}

// GetAuditLogOptions returns available filter options for the audit log
func (h *AuditHandler) GetAuditLogOptions(w http.ResponseWriter, r *http.Request) {
	// Get unique actions
	var actions []string
	h.DB.Model(&models.AuditLog{}).Distinct("action").Pluck("action", &actions)

	// Get unique entity types
	var entityTypes []string
	h.DB.Model(&models.AuditLog{}).Distinct("entity_type").Pluck("entity_type", &entityTypes)

	// Get users who have audit log entries
	type UserOption struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	var users []UserOption
	h.DB.Table("audit_logs").
		Select("DISTINCT users.id, users.first_name || ' ' || users.last_name as name").
		Joins("LEFT JOIN users ON audit_logs.user_id = users.id").
		Where("users.id IS NOT NULL").
		Order("name").
		Scan(&users)

	options := map[string]interface{}{
		"actions":      actions,
		"entity_types": entityTypes,
		"users":        users,
	}

	utils.WriteSuccess(w, options)
}