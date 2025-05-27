package routes

import (
	"hrcs/backend/config"
	"hrcs/backend/handlers"
	"hrcs/backend/middleware"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SetupRoutes(r *chi.Mux, db *gorm.DB, cfg *config.Config) {
	authHandler := handlers.NewAuthHandler(db, cfg)
	userHandler := handlers.NewUserHandler(db)
	claimHandler := handlers.NewClaimHandler(db)
	adminHandler := handlers.NewAdminHandler(db)
	dashboardHandler := handlers.NewDashboardHandler(db)

	authMiddleware := middleware.AuthMiddleware(db, cfg.JWTSecret)

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/register", authHandler.Register)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Get("/profile", userHandler.GetProfile)

			r.Route("/dashboard", func(r chi.Router) {
				r.Get("/stats", dashboardHandler.GetStats)
				r.Get("/admin-stats", dashboardHandler.GetAdminStats)
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
				})
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.AdminRequired)

				r.Get("/users", userHandler.GetUsers)
				r.Put("/users/{id}/role", userHandler.UpdateUserRole)

				r.Route("/claim-types", func(r chi.Router) {
					r.Get("/", adminHandler.GetClaimTypes)
					r.Post("/", adminHandler.CreateClaimType)
					r.Route("/{id}", func(r chi.Router) {
						r.Put("/", adminHandler.UpdateClaimType)
						r.Delete("/", adminHandler.DeleteClaimType)
					})
				})

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
					r.Delete("/{id}", adminHandler.DeleteApprovalLevel)
				})
			})
		})
	})
}
