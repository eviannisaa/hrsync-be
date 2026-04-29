package repository

import (
	"context"
	"hrsync-backend/internal/db"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/model"
)

type FeedbackRepository interface {
	GetAll(ctx context.Context, params model.ListParams) ([]dto.FeedbackResponse, int, error)
	Create(ctx context.Context, req dto.CreateFeedbackRequest) (*dto.FeedbackResponse, error)
	Delete(ctx context.Context, id string) error
}

type feedbackRepository struct {
	client *db.PrismaClient
}

func NewFeedbackRepository(client *db.PrismaClient) FeedbackRepository {
	return &feedbackRepository{client: client}
}

func (r *feedbackRepository) calculateScore(req dto.CreateFeedbackRequest) int {
	// The star rating is specifically calculated from the 'Work Experience' rate
	return req.WorkEnvironment
}

func (r *feedbackRepository) GetAll(ctx context.Context, params model.ListParams) ([]dto.FeedbackResponse, int, error) {
	skip := (params.Page - 1) * params.Limit

	var filters []db.FeedbackWhereParam
	if params.Search != "" {
		// Filter by EmployeeEmail strictly to get feedback RECEIVED BY that person
		// This ensures FeedbackOverview and FeedbackList show feedback FROM OTHERS.
		filters = append(filters,
			db.Feedback.EmployeeEmail.Contains(params.Search),
		)
	}

	sortDir := db.SortOrderAsc
	if params.SortDir == "desc" {
		sortDir = db.SortOrderDesc
	}
	var orderBy []db.FeedbackOrderByParam
	switch params.SortBy {
	case "createdAt":
		orderBy = append(orderBy, db.Feedback.CreatedAt.Order(sortDir))
	default:
		orderBy = append(orderBy, db.Feedback.CreatedAt.Order(sortDir))
	}

	dbFeedback, err := r.client.Feedback.FindMany(filters...).OrderBy(orderBy...).Skip(skip).Take(params.Limit).Exec(ctx)
	if err != nil {
		return nil, 0, err
	}

	allFeedback, err := r.client.Feedback.FindMany(filters...).Exec(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.FeedbackResponse, 0, len(dbFeedback))
	for _, du := range dbFeedback {
		createdAt := du.CreatedAt
		updatedAt := du.UpdatedAt
		createdBy, _ := du.CreatedBy()
		responses = append(responses, dto.FeedbackResponse{
			ID:                         du.ID,
			Email:                      du.Email,
			EmployeeName:               du.EmployeeName,
			EmployeeEmail:              du.EmployeeEmail,
			EmployeeDepartment:         du.EmployeeDepartment,
			Month:                      du.Month,
			IsAnonymouse:               du.IsAnonymouse,
			PositiveExperience:         du.PositiveExperience,
			Suggestion:                 du.Suggestion,
			WorkEnvironment:            du.WorkEnvironment,
			WorkQualityReliability:     du.WorkQualityReliability,
			CollaborationCommunication: du.CollaborationCommunication,
			WorkLifeBalance:            du.WorkLifeBalance,
			CriticalThinking:           du.CriticalThinking,
			OverallSatisfaction:        du.OverallSatisfaction,
			Score:                      float64(du.Score),
			CreatedAt:                  &createdAt,
			UpdatedAt:                  &updatedAt,
			CreatedBy:                  string(createdBy),
		})
	}

	return responses, len(allFeedback), nil
}

func (r *feedbackRepository) Create(ctx context.Context, req dto.CreateFeedbackRequest) (*dto.FeedbackResponse, error) {
	score := r.calculateScore(req)

	du, err := r.client.Feedback.CreateOne(
		db.Feedback.Email.Set(req.Email),
		db.Feedback.EmployeeName.Set(req.EmployeeName),
		db.Feedback.EmployeeEmail.Set(req.EmployeeEmail),
		db.Feedback.EmployeeDepartment.Set(req.EmployeeDepartment),
		db.Feedback.Month.Set(req.Month),
		db.Feedback.PositiveExperience.Set(req.PositiveExperience),
		db.Feedback.WorkEnvironment.Set(req.WorkEnvironment),
		db.Feedback.WorkQualityReliability.Set(req.WorkQualityReliability),
		db.Feedback.CollaborationCommunication.Set(req.CollaborationCommunication),
		db.Feedback.WorkLifeBalance.Set(req.WorkLifeBalance),
		db.Feedback.CriticalThinking.Set(req.CriticalThinking),
		db.Feedback.OverallSatisfaction.Set(req.OverallSatisfaction),
		db.Feedback.Score.Set(score),
		db.Feedback.IsAnonymouse.Set(req.IsAnonymouse),
		db.Feedback.Suggestion.Set(req.Suggestion),
		db.Feedback.CreatedBy.Set(req.CreatedBy),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	createdAt := du.CreatedAt
	updatedAt := du.UpdatedAt
	createdBy, _ := du.CreatedBy()
	return &dto.FeedbackResponse{
		ID:                         du.ID,
		Email:                      du.Email,
		EmployeeName:               du.EmployeeName,
		EmployeeEmail:              du.EmployeeEmail,
		EmployeeDepartment:         du.EmployeeDepartment,
		Month:                      du.Month,
		IsAnonymouse:               du.IsAnonymouse,
		PositiveExperience:         du.PositiveExperience,
		Suggestion:                 du.Suggestion,
		WorkEnvironment:            du.WorkEnvironment,
		WorkQualityReliability:     du.WorkQualityReliability,
		CollaborationCommunication: du.CollaborationCommunication,
		WorkLifeBalance:            du.WorkLifeBalance,
		CriticalThinking:           du.CriticalThinking,
		OverallSatisfaction:        du.OverallSatisfaction,
		Score:                      float64(du.Score),
		CreatedAt:                  &createdAt,
		UpdatedAt:                  &updatedAt,
		CreatedBy:                  string(createdBy),
	}, nil
}

func (r *feedbackRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Feedback.FindUnique(
		db.Feedback.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}
