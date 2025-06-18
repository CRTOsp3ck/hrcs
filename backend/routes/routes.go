package routes

import (
	"hrcs/backend/config"
	"hrcs/backend/handlers"
	"hrcs/backend/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SetupRoutes(r *chi.Mux, db *gorm.DB, cfg *config.Config) {
	authHandler := handlers.NewAuthHandler(db, cfg)
	userHandler := handlers.NewUserHandler(db)
	claimHandler := handlers.NewClaimHandler(db)
	adminHandler := handlers.NewAdminHandler(db)
	adminEnhanced := handlers.NewAdminEnhancedHandler(db)
	dashboardHandler := handlers.NewDashboardHandler(db)
	balanceHandler := handlers.NewBalanceHandler(db)
	auditHandler := handlers.NewAuditHandler(db)

	authMiddleware := middleware.AuthMiddleware(db, cfg.JWTSecret)

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/register", authHandler.Register)

		// Test endpoint
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte(`{"status":"ok"}`))
		})

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Get("/profile", userHandler.GetProfile)

			r.Route("/dashboard", func(r chi.Router) {
				r.Get("/stats", dashboardHandler.GetStats)
				r.Get("/admin-stats", dashboardHandler.GetAdminStats)
			})

			// Claim types for regular users (read-only)
			r.Get("/claim-types", claimHandler.GetClaimTypes)

			// Balance routes for regular users
			r.Route("/balances", func(r chi.Router) {
				r.Get("/", balanceHandler.GetUserBalances)
				r.Get("/claim-type/{id}", balanceHandler.GetUserBalance)
				r.Post("/check", balanceHandler.CheckClaimAmount)
			})

			r.Route("/claims", func(r chi.Router) {
				r.Get("/", claimHandler.GetClaims)
				r.Post("/", claimHandler.CreateClaim)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", claimHandler.GetClaim)
					r.Put("/", claimHandler.UpdateClaim)
					r.Delete("/", claimHandler.CancelClaim)
					r.Post("/submit", claimHandler.SubmitClaim)
					r.Post("/approve", claimHandler.ApproveClaim)
					// NEW: Enhanced workflow endpoints (Phase 3)
					r.Get("/workflow", claimHandler.GetClaimWorkflow)
					r.Put("/workflow/{stepId}", claimHandler.UpdateClaimWorkflowStep)
				})
			})

			// Admin routes with /admin prefix
			r.Group(func(r chi.Router) {
				r.Use(middleware.AdminRequired)

				r.Route("/admin", func(r chi.Router) {
					// Claims management
					r.Route("/claims", func(r chi.Router) {
						r.Get("/", adminEnhanced.GetAllClaims)
						r.Post("/{id}/approve", adminEnhanced.AdminApproveClaim)
						r.Post("/{id}/reject", adminEnhanced.AdminRejectClaim)
						r.Put("/{id}/status", adminEnhanced.UpdateClaimStatus)
					})

					// Users management
					r.Route("/users", func(r chi.Router) {
						r.Get("/", adminEnhanced.GetAdminUsers)
						r.Post("/", adminEnhanced.CreateAdminUser)
						r.Put("/{id}", adminEnhanced.UpdateAdminUser)
						r.Delete("/{id}", adminEnhanced.DeleteAdminUser)
						// NEW: User claim overrides and balance details
						r.Post("/{id}/claim-overrides", adminHandler.SetUserClaimOverrides)
						r.Get("/{id}/balance-details", balanceHandler.GetUserBalanceDetails)
						// NEW: User detail view
						r.Get("/{id}/details", adminHandler.GetUserDetails)
					})

					// Groups management
					r.Route("/groups", func(r chi.Router) {
						r.Get("/", adminEnhanced.GetEnhancedGroups)
						r.Post("/", adminEnhanced.CreateEnhancedGroup)
						r.Put("/{id}", adminEnhanced.UpdateEnhancedGroup)
						r.Delete("/{id}", adminEnhanced.DeleteEnhancedGroup)
						// NEW: Group claim permissions
						r.Post("/{id}/claim-permissions", adminHandler.SetUserGroupClaimPermissions)
						// NEW: Group detail view
						r.Get("/{id}/details", adminHandler.GetUserGroupDetails)
					})

					// Claim Types management
					r.Route("/claim-types", func(r chi.Router) {
						r.Get("/", adminEnhanced.GetEnhancedClaimTypes)
						r.Post("/", adminEnhanced.CreateEnhancedClaimType)
						r.Put("/{id}", adminEnhanced.UpdateEnhancedClaimType)
						r.Delete("/{id}", adminHandler.DeleteClaimType)
						// NEW: Claim type limits
						r.Put("/{id}/limits", adminHandler.UpdateClaimTypeWithLimits)
						// NEW: Claim type detail view
						r.Get("/{id}/details", adminHandler.GetClaimTypeDetails)
					})

					// Approval Levels management
					r.Route("/approval-levels", func(r chi.Router) {
						r.Get("/", adminEnhanced.GetEnhancedApprovalLevels)
						r.Get("/by-group", adminEnhanced.GetApprovalLevelsByGroup)
						r.Post("/", adminEnhanced.CreateEnhancedApprovalLevel)
						r.Put("/{id}", adminEnhanced.UpdateEnhancedApprovalLevel)
						r.Delete("/{id}", adminEnhanced.DeleteEnhancedApprovalLevel)
						r.Put("/order", adminEnhanced.UpdateApprovalLevelOrder)
					})

					// NEW: Balance management
					r.Route("/balances", func(r chi.Router) {
						r.Post("/adjust", balanceHandler.AdminUpdateBalance)
					})

					// NEW: Audit log management (Phase 3)
					r.Route("/audit-log", func(r chi.Router) {
						r.Get("/", auditHandler.GetAuditLogs)
						r.Get("/stats", auditHandler.GetAuditLogStats)
						r.Get("/options", auditHandler.GetAuditLogOptions)
						r.Get("/{id}", auditHandler.GetAuditLogDetails)
					})

					// NEW: Reports management (Phase 3)
					r.Route("/reports", func(r chi.Router) {
						// Placeholder endpoints for reports
						r.Get("/", func(w http.ResponseWriter, r *http.Request) {
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusOK)
							w.Write([]byte(`{"data": [], "message": "Reports feature coming soon"}`))
						})
						r.Post("/{type}/generate", func(w http.ResponseWriter, r *http.Request) {
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusOK)
							w.Write([]byte(`{"message": "Report generation feature coming soon"}`))
						})
					})

					// NEW: Integrations management (Phase 3)
					r.Route("/integrations", func(r chi.Router) {
						// Placeholder endpoints for integrations
						r.Get("/", func(w http.ResponseWriter, r *http.Request) {
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusOK)
							w.Write([]byte(`{"data": {"microsoft": {"configured": false, "features": ["sso", "teams", "ad_sync"]}}, "message": "Integrations configuration available"}`))
						})
						r.Put("/microsoft", func(w http.ResponseWriter, r *http.Request) {
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusOK)
							w.Write([]byte(`{"message": "Microsoft integration configuration coming soon"}`))
						})
					})
				})

				// Legacy routes (keeping for backward compatibility)
				r.Get("/users", userHandler.GetUsers)
				r.Put("/users/{id}/role", userHandler.UpdateUserRole)

				r.Route("/user-groups", func(r chi.Router) {
					r.Get("/", adminHandler.GetUserGroups)
					r.Post("/", adminHandler.CreateUserGroup)
					r.Route("/{id}", func(r chi.Router) {
						r.Put("/", adminHandler.UpdateUserGroup)
						r.Delete("/", adminHandler.DeleteUserGroup)
					})
				})

				r.Route("/approval-levels", func(r chi.Router) {
					r.Get("/", adminHandler.GetApprovalLevels)
					r.Post("/", adminHandler.CreateApprovalLevel)
					r.Put("/{id}", adminHandler.UpdateApprovalLevel)
					r.Delete("/{id}", adminHandler.DeleteApprovalLevel)
				})
			})
		})
	})
}
