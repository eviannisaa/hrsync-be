package service

import (
	"context"
	"encoding/json"
	"fmt"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/repository"
	"log"
	"net/http"
	"time"
)

type HolidayService interface {
	GetHolidays(ctx context.Context) ([]dto.HolidayResponse, error)
	SyncHolidays(ctx context.Context, year int) error
}

type holidayService struct {
	repo repository.HolidayRepository
}

func NewHolidayService(repo repository.HolidayRepository) HolidayService {
	return &holidayService{repo: repo}
}

func (s *holidayService) GetHolidays(ctx context.Context) ([]dto.HolidayResponse, error) {
	return s.repo.GetAll(ctx)
}

func (s *holidayService) SyncHolidays(ctx context.Context, year int) error {
	url := fmt.Sprintf("https://tanggalmerah.upset.dev/api/holidays?year=%d", year)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[HolidayService] Error fetching holidays for %d: %v", year, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[HolidayService] Error status code for %d: %d", year, resp.StatusCode)
		return fmt.Errorf("failed to fetch holidays: status %d", resp.StatusCode)
	}

	var response dto.SyncHolidayResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("[HolidayService] Error decoding response for %d: %v", year, err)
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
