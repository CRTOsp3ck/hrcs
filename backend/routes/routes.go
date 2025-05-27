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

			r.Route("/claims", func(r chi.Router) {
				r.Get("/", claimHandler.GetClaims)
				r.Post("/", claimHandler.CreateClaim)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", claimHandler.GetClaim)
					r.Put("/", claimHandler.UpdateClaim)
					r.Delete("/", claimHandler.CancelClaim)
					r.Post("/submit", claimHandler.SubmitClaim)
					r.Post("/approve", claimHandler.ApproveClaim)
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
					})

					// Groups management
					r.Route("/groups", func(r chi.Router) {
						r.Get("/", adminEnhanced.GetEnhancedGroups)
						r.Post("/", adminEnhanced.CreateEnhancedGroup)
						r.Put("/{id}", adminEnhanced.UpdateEnhancedGroup)
						r.Delete("/{id}", adminEnhanced.DeleteEnhancedGroup)
					})

					// Claim Types management
					r.Route("/claim-types", func(r chi.Router) {
						r.Get("/", adminEnhanced.GetEnhancedClaimTypes)
						r.Post("/", adminEnhanced.CreateEnhancedClaimType)
						r.Put("/{id}", adminEnhanced.UpdateEnhancedClaimType)
						r.Delete("/{id}", adminHandler.DeleteClaimType)
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
