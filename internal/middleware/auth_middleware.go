package middleware

import (
	"context"
	"hrsync-backend/internal/utils"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "userID"
	ContextKeyEmail  contextKey = "email"
	ContextKeyRole   contextKey = "role"
)

// AuthMiddleware memvalidasi JWT token dari header Authorization: Bearer <token>.
// Jika valid, menyimpan claims ke dalam context request.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := ""

		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				token = parts[1]
			}
		}

		// Jika tidak ada di header, cek cookie
		if token == "" {
			cookie, err := r.Cookie("hrsync_token")
			if err == nil {
				token = cookie.Value
			}
		}

		if token == "" {
			utils.SendError(w, "authorization required", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			utils.SendError(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Simpan claims ke context
		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyEmail, claims.Email)
		ctx = context.WithValue(ctx, ContextKeyRole, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleMiddleware memastikan user yang login memiliki salah satu role yang diizinkan.
// Harus digunakan setelah AuthMiddleware.
func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(ContextKeyRole).(string)
			if !ok {
				utils.SendError(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if _, allowed := roleSet[role]; !allowed {
				utils.SendError(w, "forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
