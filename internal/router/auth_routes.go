package router

import (
	"hrsync-backend/internal/handler"
	"hrsync-backend/internal/middleware"
	"hrsync-backend/internal/utils"
	"net/http"
)

func RegisterAuthRoutes(mux *http.ServeMux, authHandler *handler.AuthHandler) {
	// Public routes — tidak butuh auth
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)

	mux.Handle("GET /api/auth/me", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(middleware.ContextKeyUserID).(string)
		email, _ := r.Context().Value(middleware.ContextKeyEmail).(string)
		role, _ := r.Context().Value(middleware.ContextKeyRole).(string)

		// Frontend ekspektasi response.data.user
		// token tidak perlu dikirim ulang di /me
		data := map[string]interface{}{
			"user": map[string]string{
				"id":    userID,
				"email": email,
				"role":  role,
			},
		}

		utils.SendSuccess(w, "authenticated", data, http.StatusOK)
	})))
}
