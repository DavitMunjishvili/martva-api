package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"driving-license-city-api/internal/models"
)

// SAService is a service for fetching data from the sa.gov.ge API.
type SAService struct {
	Client *http.Client
}

// NewSAService creates a new SAService.
func NewSAService(client *http.Client) *SAService {
	return &SAService{Client: client}
}

// FetchDatesForCenter fetches available dates for a given center.
func (s *SAService) FetchDatesForCenter(centerID int, centerName string, ch chan<- models.CenterResult, wg *sync.WaitGroup) {
	defer wg.Done()

	result, err := s.FetchDatesForCenterSync(centerID, centerName)
	if err != nil {
		result.Error = err.Error()
	}

	ch <- result
}

// FetchDatesForCenterSync fetches available dates for a given center synchronously.
func (s *SAService) FetchDatesForCenterSync(centerID int, centerName string) (models.CenterResult, error) {
	url := fmt.Sprintf("https://api-my.sa.gov.ge/api/v1/DrivingLicensePracticalExams2/DrivingLicenseExamsDates2?CategoryCode=4&CenterId=%d", centerID)
	result := models.CenterResult{CenterName: centerName, CenterID: centerID}

	resp, err := s.Client.Get(url)
	if err != nil {
		return result, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("API returned non-200 status: %s", resp.Status)
	}

	var datesResponse []models.AvailableDate
	if err := json.NewDecoder(resp.Body).Decode(&datesResponse); err != nil {
		return result, fmt.Errorf("error decoding JSON: %v", err)
	}

	result.Dates = datesResponse
	return result, nil
}
