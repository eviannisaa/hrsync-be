package service

import (
	"context"
	"fmt"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/repository"
	"hrsync-backend/pkg/external"
	"log"
	"time"
)

type HolidayService interface {
	GetHolidays(ctx context.Context) ([]dto.HolidayResponse, error)
	SyncHolidays(ctx context.Context, year int) error
}

type holidayService struct {
	repo repository.HolidayRepository
	api  external.HolidayAPI
}

func NewHolidayService(repo repository.HolidayRepository, api external.HolidayAPI) HolidayService {
	return &holidayService{
		repo: repo,
		api:  api,
	}
}

func (s *holidayService) GetHolidays(ctx context.Context) ([]dto.HolidayResponse, error) {
	return s.repo.GetAll(ctx)
}

func (s *holidayService) SyncHolidays(ctx context.Context, year int) error {
	response, err := s.api.FetchHolidays(year)
	if err != nil {
		log.Printf("[HolidayService] Error fetching holidays for %d: %v", year, err)
		return err
	}

	if !response.Success {
		log.Printf("[HolidayService] API success false for %d", year)
		return fmt.Errorf("API success false")
	}

	log.Printf("[HolidayService] Syncing %d holidays for year %d", len(response.Data), year)

	for _, item := range response.Data {
		date, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			log.Printf("[HolidayService] Error parsing date %s: %v", item.Date, err)
			continue
		}

		isCollective := item.Type == "leave"
		_, err = s.repo.Upsert(ctx, item.Name, date, isCollective)
		if err != nil {
			log.Printf("[HolidayService] Error upserting holiday %s: %v", item.Name, err)
		}
	}

	return nil
}
