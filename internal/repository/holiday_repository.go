package repository

import (
	"context"
	"hrsync-backend/internal/db"
	"hrsync-backend/internal/dto"
	"time"
)

type HolidayRepository interface {
	GetAll(ctx context.Context) ([]dto.HolidayResponse, error)
	Upsert(ctx context.Context, name string, date time.Time, isCollective bool) (*dto.HolidayResponse, error)
}

type holidayRepository struct {
	client *db.PrismaClient
}

func NewHolidayRepository(client *db.PrismaClient) HolidayRepository {
	return &holidayRepository{client: client}
}

func (r *holidayRepository) GetAll(ctx context.Context) ([]dto.HolidayResponse, error) {
	dbHolidays, err := r.client.Holiday.FindMany().OrderBy(
		db.Holiday.Date.Order(db.SortOrderAsc),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.HolidayResponse, 0, len(dbHolidays))
	for _, du := range dbHolidays {
		responses = append(responses, dto.HolidayResponse{
			InnerHoliday: du.InnerHoliday,
			CreatedAt:    &du.CreatedAt,
			UpdatedAt:    &du.UpdatedAt,
		})
	}

	return responses, nil
}

func (r *holidayRepository) Upsert(ctx context.Context, name string, date time.Time, isCollective bool) (*dto.HolidayResponse, error) {
	du, err := r.client.Holiday.UpsertOne(
		db.Holiday.Date.Equals(date),
	).Create(
		db.Holiday.Name.Set(name),
		db.Holiday.Date.Set(date),
		db.Holiday.IsCollective.Set(isCollective),
	).Update(
		db.Holiday.Name.Set(name),
		db.Holiday.IsCollective.Set(isCollective),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &dto.HolidayResponse{
		InnerHoliday: du.InnerHoliday,
		CreatedAt:    &du.CreatedAt,
		UpdatedAt:    &du.UpdatedAt,
	}, nil
}
