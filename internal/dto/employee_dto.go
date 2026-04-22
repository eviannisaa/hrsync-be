package dto

import (
	"hrsync-backend/internal/db"
	"time"
)

type EmployeeResponse struct {
	db.InnerEmployee
	Status    string     `json:"status"`
	CreatedBy *string    `json:"createdBy,omitempty"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type CreateEmployeeRequest struct {
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Department string    `json:"department"`
	Position   string    `json:"position"`
	JoinDate   time.Time `json:"joinDate"`
	Location   string    `json:"location"`
	Status     string    `json:"status"`
	CreatedBy  string    `json:"createdBy,omitempty"`
	UpdatedBy  string    `json:"updatedBy,omitempty"`
}

type UpdateEmployeeRequest struct {
	Name       *string    `json:"name,omitempty"`
	Email      *string    `json:"email,omitempty"`
	Phone      *string    `json:"phone,omitempty"`
	Department *string    `json:"department,omitempty"`
	Position   *string    `json:"position,omitempty"`
	JoinDate   *time.Time `json:"joinDate,omitempty"`
	IsActive   *bool      `json:"isActive,omitempty"`
	Status     *string    `json:"status,omitempty"`
	Latitude   *float64   `json:"latitude,omitempty"`
	Longitude  *float64   `json:"longitude,omitempty"`
	Location   *string    `json:"location,omitempty"`
	UpdatedBy  *string    `json:"updatedBy,omitempty"`
}

type EmployeeOrganizationResponse struct {
	ID                string    `json:"id"`
	OrganizationImage string    `json:"organizationImage"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type UpdateEmployeeOrganizationRequest struct {
	OrganizationImage string `json:"organizationImage"`
}
