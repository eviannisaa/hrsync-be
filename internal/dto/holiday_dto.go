package dto

import (
	"hrsync-backend/internal/db"
	"time"
)

type HolidayResponse struct {
	db.InnerHoliday
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type SyncHolidayResponse struct {
	Success bool              `json:"success"`
	Data    []SyncHolidayItem `json:"data"`
}

type SyncHolidayItem struct {
	Date string `json:"date"`
	Name string `json:"name"`
	Type string `json:"type"`
}
