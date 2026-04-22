package seeding

import (
	"context"
	"hrsync-backend/internal/db"
	"log"
)

func SeedFeedbacks(ctx context.Context, client *db.PrismaClient) {
	feedbacks := []struct {
		Email                      string
		EmployeeName               string
		EmployeeEmail              string
		EmployeeDepartment         string
		Month                      string
		IsAnonymouse               bool
		PositiveExperience         string
		Suggestion                 string
		WorkEnvironment            int
		WorkQualityReliability     int
		CollaborationCommunication int
		WorkLifeBalance            int
		CriticalThinking           int
		OverallSatisfaction        int
		Score                      int
	}{
		{
			Email:                      "john.doe@example.com",
			EmployeeName:               "John Doe",
			EmployeeEmail:              "john.doe@example.com",
			EmployeeDepartment:         "Engineering",
			Month:                      "April 2026",
			IsAnonymouse:               false,
			PositiveExperience:         "The new project structure is very helpful.",
			Suggestion:                 "None so far.",
			WorkEnvironment:            5,
			WorkQualityReliability:     5,
			CollaborationCommunication: 4,
			WorkLifeBalance:            4,
			CriticalThinking:           5,
			OverallSatisfaction:        5,
			Score:                      5,
		},
		{
			Email:                      "jane.smith@example.com",
			EmployeeName:               "Jane Smith",
			EmployeeEmail:              "jane.smith@example.com",
			EmployeeDepartment:         "HR",
			Month:                      "April 2026",
			IsAnonymouse:               true,
			PositiveExperience:         "I like the teamwork in our department.",
			Suggestion:                 "Improve the kitchen facilities.",
			WorkEnvironment:            4,
			WorkQualityReliability:     4,
			CollaborationCommunication: 5,
			WorkLifeBalance:            3,
			CriticalThinking:           4,
			OverallSatisfaction:        4,
			Score:                      4,
		},
	}

	for _, f := range feedbacks {
		// Cek apakah sudah ada (berdasarkan email dan month untuk menghindari duplikasi data seed)
		existing, _ := client.Feedback.FindFirst(
			db.Feedback.Email.Equals(f.Email),
			db.Feedback.Month.Equals(f.Month),
		).Exec(ctx)
		if existing != nil {
			log.Printf("Feedback record already exists for %s, skipping", f.Email)
			continue
		}

		_, err := client.Feedback.CreateOne(
			db.Feedback.Email.Set(f.Email),
			db.Feedback.EmployeeName.Set(f.EmployeeName),
			db.Feedback.EmployeeEmail.Set(f.EmployeeEmail),
			db.Feedback.EmployeeDepartment.Set(f.EmployeeDepartment),
			db.Feedback.Month.Set(f.Month),
			db.Feedback.PositiveExperience.Set(f.PositiveExperience),
			db.Feedback.WorkEnvironment.Set(f.WorkEnvironment),
			db.Feedback.WorkQualityReliability.Set(f.WorkQualityReliability),
			db.Feedback.CollaborationCommunication.Set(f.CollaborationCommunication),
			db.Feedback.WorkLifeBalance.Set(f.WorkLifeBalance),
			db.Feedback.CriticalThinking.Set(f.CriticalThinking),
			db.Feedback.OverallSatisfaction.Set(f.OverallSatisfaction),
			db.Feedback.Score.Set(f.Score),
			db.Feedback.IsAnonymouse.Set(f.IsAnonymouse),
			db.Feedback.Suggestion.Set(f.Suggestion),
		).Exec(ctx)
		if err != nil {
			log.Printf("failed to create feedback for %s: %v", f.Email, err)
		} else {
			log.Printf("Created feedback record for email: %s", f.Email)
		}
	}
}
