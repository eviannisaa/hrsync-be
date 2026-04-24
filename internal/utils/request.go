package utils

import (
	"hrsync-backend/internal/model"
	"net/http"
	"strconv"
)

// GetListParams extracts common list parameters from the request query
func GetListParams(r *http.Request) model.ListParams {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	email, _ := r.Context().Value(model.ContextKeyEmail).(string)
	role, _ := r.Context().Value(model.ContextKeyRole).(string)

	return model.ListParams{
		Page:    page,
		Limit:   limit,
		Search:  r.URL.Query().Get("search"),
		SortBy:  r.URL.Query().Get("sortBy"),
		SortDir: r.URL.Query().Get("sortDir"),
		Email:   email,
		Role:    role,
	}
}
