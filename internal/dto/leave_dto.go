package dto

import (
	"hrsync-backend/internal/db"
	"time"
)

type LeaveResponse struct {
	db.InnerLeave
	LeaveType    string     `json:"leaveType"`
	EmployeeName string     `json:"employeeName"`
	Department   string     `json:"department"`
	CreatedBy    *string    `json:"createdBy,omitempty"`
	UpdatedBy    *string    `json:"updatedBy,omitempty"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
}

type CreateLeaveRequest struct {
	Email       string    `json:"email"`
	LeaveTypeId string    `json:"leaveTypeId"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Reason      string    `json:"reason"`
	CreatedBy   string    `json:"createdBy,omitempty"`
	UpdatedBy   string    `json:"updatedBy,omitempty"`
}

type UpdateLeaveRequest struct {
	Email       *string    `json:"email"`
	LeaveTypeId *string    `json:"leaveTypeId"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Reason      *string    `json:"reason"`
	UpdatedBy   string     `json:"updatedBy,omitempty"`
}

type LeaveStatusSummary struct {
	Status string `json:"status"`
	Total  int    `json:"total"`
}

type LeaveSummaryResponse struct {
	Summary []LeaveStatusSummary `json:"summary"`
}
