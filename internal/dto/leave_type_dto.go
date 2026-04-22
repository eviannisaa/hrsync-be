package dto

import (
	"hrsync-backend/internal/db"
	"time"
)

type LeaveTypeResponse struct {
	db.InnerLeaveType
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type LeaveBalanceResponse struct {
	Leave     string `json:"leave"`
	Total     int    `json:"total"`
	Used      int    `json:"used"`
	Remaining int    `json:"remaining"`
}

type CreateLeaveTypeRequest struct {
	Name        string `json:"name"`
	DefaultDays int    `json:"defaultDays"`
}
