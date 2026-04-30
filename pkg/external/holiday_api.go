package external

import (
	"encoding/json"
	"fmt"
	"hrsync-backend/internal/dto"
	"net/http"
)

type HolidayAPI interface {
	FetchHolidays(year int) (*dto.SyncHolidayResponse, error)
}

type holidayAPI struct {
	baseURL string
}

func NewHolidayAPI() HolidayAPI {
	return &holidayAPI{
		baseURL: "https://tanggalmerah.upset.dev/api/holidays",
	}
}

func (a *holidayAPI) FetchHolidays(year int) (*dto.SyncHolidayResponse, error) {
	url := fmt.Sprintf("%s?year=%d", a.baseURL, year)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API error: status %d", resp.StatusCode)
	}

	var result dto.SyncHolidayResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
