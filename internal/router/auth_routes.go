package router

import (
	"hrsync-backend/internal/handler"
	"hrsync-backend/internal/middleware"
	"hrsync-backend/internal/model"
	"hrsync-backend/internal/repository"
	"hrsync-backend/internal/utils"
	"net/http"
)

func RegisterAuthRoutes(mux *http.ServeMux, authHandler *handler.AuthHandler, employeeRepo repository.EmployeeRepository) {
	// Public routes — tidak butuh auth
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)

	mux.Handle("GET /api/auth/me", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(model.ContextKeyUserID).(string)
		email, _ := r.Context().Value(model.ContextKeyEmail).(string)
		role, _ := r.Context().Value(model.ContextKeyRole).(string)

		// Fetch employee info for department info
		employee, _ := employeeRepo.GetByEmail(r.Context(), email)

		// Frontend ekspektasi response.data.user
		// token tidak perlu dikirim ulang di /me
		data := map[string]interface{}{
			"user": map[string]interface{}{
				"id":       userID,
				"email":    email,
				"role":     role,
				"employee": employee,
			},
		}

		utils.SendSuccess(w, "authenticated", data, http.StatusOK)
	})))
}
