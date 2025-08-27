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

	url := fmt.Sprintf("https://api-my.sa.gov.ge/api/v1/DrivingLicensePracticalExams2/DrivingLicenseExamsDates2?CategoryCode=4&CenterId=%d", centerID)
	result := models.CenterResult{CenterName: centerName}

	resp, err := s.Client.Get(url)
	if err != nil {
		result.Error = fmt.Sprintf("Error fetching data: %v", err)
		ch <- result
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Error = fmt.Sprintf("API returned non-200 status: %s", resp.Status)
		ch <- result
		return
	}

	var datesResponse []models.AvailableDate
	if err := json.NewDecoder(resp.Body).Decode(&datesResponse); err != nil {
		result.Error = fmt.Sprintf("Error decoding JSON: %v", err)
		ch <- result
		return
	}

	result.Dates = datesResponse
	result.CenterID = centerID
	ch <- result
}
