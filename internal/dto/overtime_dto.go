package dto

import (
	"hrsync-backend/internal/db"
	"time"
)

type OvertimeResponse struct {
	db.InnerOvertime
	EmployeeName string     `json:"employeeName"`
	Department   string     `json:"department"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
}

type CreateOvertimeRequest struct {
	Email       string    `json:"email"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	StartTime   string    `json:"startTime"`
	EndTime     string    `json:"endTime"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
}

type UpdateOvertimeRequest struct {
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	StartTime   *string    `json:"startTime"`
	EndTime     *string    `json:"endTime"`
	Type        *string    `json:"type"`
	Description *string    `json:"description"`
}
