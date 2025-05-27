package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"hrcs/backend/models"
	"hrcs/backend/utils"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(db *gorm.DB, jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				utils.WriteError(w, http.StatusUnauthorized, "Bearer token required")
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				utils.WriteError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				utils.WriteError(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}

			userID, ok := claims["user_id"].(float64)
			if !ok {
				utils.WriteError(w, http.StatusUnauthorized, "Invalid user ID in token")
				return
			}

			var user models.User
			if err := db.Preload("UserGroup").First(&user, uint(userID)).Error; err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "User not found")
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		if user == nil || user.Role != models.RoleAdmin {
			utils.WriteError(w, http.StatusForbidden, "Admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserFromContext(ctx context.Context) *models.User {
	user, ok := ctx.Value(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}