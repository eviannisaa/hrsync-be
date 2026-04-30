package dto

import (
	"hrsync-backend/internal/db"
	"time"
)

type FeedbackResponse struct {
	db.InnerFeedback
	Score     float64    `json:"score"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	CreatedBy string     `json:"createdBy"`
}

type CreateFeedbackRequest struct {
	Email                      string `json:"email"`
	CreatedBy                  string `json:"createdBy"`
	EmployeeName               string `json:"employeeName"`
	EmployeeEmail              string `json:"employeeEmail"`
	EmployeeDepartment         string `json:"employeeDepartment"`
	Month                      string `json:"month"`
	IsAnonymouse               bool   `json:"isAnonymouse"`
	PositiveExperience         string `json:"positiveExperience"`
	Suggestion                 string `json:"suggestion"`
	WorkEnvironment            int    `json:"workEnvironment"`
	WorkQualityReliability     int    `json:"workQualityReliability"`
	CollaborationCommunication int    `json:"collaborationCommunication"`
	WorkLifeBalance            int    `json:"workLifeBalance"`
	CriticalThinking           int    `json:"criticalThinking"`
	OverallSatisfaction        int    `json:"overallSatisfaction"`
}
