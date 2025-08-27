package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"driving-license-city-api/internal/config"
	"driving-license-city-api/internal/models"
	"driving-license-city-api/internal/service"
)

// API holds the dependencies for the API handlers.
type API struct {
	SAService *service.SAService
}

// NewAPI creates a new API.
func NewAPI(saService *service.SAService) *API {
	return &API{SAService: saService}
}

func (a *API) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	health := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	if err := json.NewEncoder(w).Encode(health); err != nil {
		log.Printf("Error encoding health check response: %v", err)
		http.Error(w, "Failed to create health check response", http.StatusInternalServerError)
	}
}

// AvailableDatesHandler handles the request for available dates.
func (a *API) AvailableDatesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var wg sync.WaitGroup
	resultsChannel := make(chan models.CenterResult, len(config.Centers))
	for id, name := range config.Centers {
		wg.Add(1)
		go a.SAService.FetchDatesForCenter(id, name, resultsChannel, &wg)
	}
	go func() {
		wg.Wait()
		close(resultsChannel)
	}()
	finalResponse := make(map[string]models.CenterResult)
	for result := range resultsChannel {
		finalResponse[result.CenterName] = result
	}
	if err := json.NewEncoder(w).Encode(finalResponse); err != nil {
		log.Printf("Error encoding final JSON response: %v", err)
		http.Error(w, "Failed to create response", http.StatusInternalServerError)
	}
}

// AvailableHoursHandler handles the request for available hours.
func (a *API) AvailableHoursHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	centerId := r.URL.Query().Get("centerId")
	examDate := r.URL.Query().Get("examDate")

	if centerId == "" || examDate == "" {
		http.Error(w, "Missing required query parameters: 'centerId' and 'examDate'", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://api-my.sa.gov.ge/api/v1/DrivingLicensePracticalExams2/DrivingLicenseExamsDateFrames2?CategoryCode=4&CenterId=%s&ExamDate=%s", centerId, examDate)

	resp, err := a.SAService.Client.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch data from external API: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		return
	}

	io.Copy(w, resp.Body)
}

func (a *API) CityInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	centerIdStr := r.URL.Query().Get("centerId")
	centerId, err := strconv.Atoi(centerIdStr)
	if err != nil {
		http.Error(w, "Invalid centerId", http.StatusBadRequest)
		return
	}

	centerName, ok := config.Centers[centerId]
	if !ok {
		http.Error(w, "Center not found", http.StatusNotFound)
		return
	}

	result, err := a.SAService.FetchDatesForCenterSync(centerId, centerName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding city info: %v", err)
		http.Error(w, "Failed to create city info response", http.StatusInternalServerError)
	}
}
