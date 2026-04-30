package handler

import (
	"encoding/json"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/service"
	"hrsync-backend/internal/utils"
	"net/http"
	"time"
)

type AuthHandler struct {
	srv service.AuthService
}

func NewAuthHandler(srv service.AuthService) *AuthHandler {
	return &AuthHandler{srv: srv}
}

// Register godoc
// @Summary Register akun baru
// @Description Membuat akun User baru dengan role ADMIN atau EMPLOYEE
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register payload"
// @Success 201 {object} utils.APIResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.SendError(w, "email and password are required", http.StatusBadRequest)
		return
	}

	resp, err := h.srv.Register(r.Context(), req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "hrsync_token",
		Value:    resp.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set true in production
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	utils.SendSuccess(w, "Registration successful", resp, http.StatusCreated)
}

// Logout godoc
// @Summary Logout dari sistem
// @Description Menghapus session cookie
// @Tags auth
// @Produce json
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "hrsync_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	utils.SendSuccess(w, "Logout successful", nil, http.StatusOK)
}

// Login godoc
// @Summary Login ke sistem
// @Description Autentikasi user dan mendapatkan JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login payload"
// @Success 200 {object} utils.APIResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.SendError(w, "email and password are required", http.StatusBadRequest)
		return
	}

	resp, err := h.srv.Login(r.Context(), req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "hrsync_token",
		Value:    resp.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set true in production
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	utils.SendSuccess(w, "Login successful", resp, http.StatusOK)
}

// GeneratePassword godoc
// @Summary Generate password untuk employee
// @Description Menghasilkan password acak dan membuat akun User untuk employee
// @Tags auth
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} utils.APIResponse
// @Router /employees/{id}/generate-password [post]
func (h *AuthHandler) GeneratePassword(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.SendError(w, "Employee ID is required", http.StatusBadRequest)
		return
	}

	password, err := h.srv.GeneratePassword(r.Context(), id)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Password generated successfully", map[string]string{
		"password": password,
	}, http.StatusOK)
}
