package handler

import (
	"encoding/json"
	"fmt"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/service"
	"hrsync-backend/internal/utils"
	"net/http"
	"net/url"
	"os"
	"time"
)

func getGoogleAuthURL() string {
	baseURL := "https://accounts.google.com/o/oauth2/v2/auth"
	v := url.Values{}
	v.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	v.Set("redirect_uri", os.Getenv("GOOGLE_REDIRECT_URL"))
	v.Set("response_type", "code")
	v.Set("scope", "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile")
	v.Set("state", "state") // Should be random in production
	return baseURL + "?" + v.Encode()
}

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

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[AuthHandler] GoogleLogin hit")
	http.Redirect(w, r, getGoogleAuthURL(), http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[AuthHandler] GoogleCallback hit")
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.SendError(w, "Code not found", http.StatusBadRequest)
		return
	}

	// 1. Exchange code for token
	tokenRes, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"code":          {code},
		"client_id":     {os.Getenv("GOOGLE_CLIENT_ID")},
		"client_secret": {os.Getenv("GOOGLE_CLIENT_SECRET")},
		"redirect_uri":  {os.Getenv("GOOGLE_REDIRECT_URL")},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		utils.SendError(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}
	defer tokenRes.Body.Close()

	var tokenData struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(tokenRes.Body).Decode(&tokenData); err != nil {
		utils.SendError(w, "Failed to decode token", http.StatusInternalServerError)
		return
	}

	// 2. Get user info
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)
	userRes, err := client.Do(req)
	if err != nil {
		utils.SendError(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer userRes.Body.Close()

	var userinfo struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(userRes.Body).Decode(&userinfo); err != nil {
		utils.SendError(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	resp, err := h.srv.HandleGoogleAuth(r.Context(), userinfo.Email)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Set cookie and redirect
	http.SetCookie(w, &http.Cookie{
		Name:     "hrsync_token",
		Value:    resp.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	http.Redirect(w, r, "http://localhost:3000/overview", http.StatusTemporaryRedirect)
}
