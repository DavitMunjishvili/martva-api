package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type availableDate struct {
	BookingDate       string `json:"bookingDate"`
	BookingDateStatus int    `json:"bookingDateStatus"`
}

type centerResult struct {
	CenterId   int             `json:"centerId"`
	CenterName string          `json:"centerName"`
	Dates      []availableDate `json:"dates"`
	Error      string          `json:"error,omitempty"`
}

var centers = map[int]string{
	2:  "Kutaisi",
	3:  "Batumi",
	4:  "Telavi",
	5:  "Akhaltsikhe",
	6:  "Zugdidi",
	7:  "Gori",
	8:  "Poti",
	9:  "Ozurgeti",
	10: "Sachkhere",
	15: "Rustavi",
}

func fetchDatesForCenter(centerID int, centerName string, ch chan<- centerResult, wg *sync.WaitGroup) {
	defer wg.Done()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	url := fmt.Sprintf("https://api-my.sa.gov.ge/api/v1/DrivingLicensePracticalExams2/DrivingLicenseExamsDates2?CategoryCode=4&CenterId=%d", centerID)
	result := centerResult{CenterName: centerName}

	resp, err := client.Get(url)
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

	var datesResponse []availableDate
	if err := json.NewDecoder(resp.Body).Decode(&datesResponse); err != nil {
		result.Error = fmt.Sprintf("Error decoding JSON: %v", err)
		ch <- result
		return
	}

	result.Dates = datesResponse
	result.CenterId = centerID
	ch <- result
}

func availableDatesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var wg sync.WaitGroup
	resultsChannel := make(chan centerResult, len(centers))
	for id, name := range centers {
		wg.Add(1)
		go fetchDatesForCenter(id, name, resultsChannel, &wg)
	}
	go func() {
		wg.Wait()
		close(resultsChannel)
	}()
	finalResponse := make(map[string]centerResult)
	for result := range resultsChannel {
		finalResponse[result.CenterName] = result
	}
	if err := json.NewEncoder(w).Encode(finalResponse); err != nil {
		log.Printf("Error encoding final JSON response: %v", err)
		http.Error(w, "Failed to create response", http.StatusInternalServerError)
	}
}

func availableHoursHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	centerId := r.URL.Query().Get("centerId")
	examDate := r.URL.Query().Get("examDate")

	if centerId == "" || examDate == "" {
		http.Error(w, "Missing required query parameters: 'centerId' and 'examDate'", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://api-my.sa.gov.ge/api/v1/DrivingLicensePracticalExams2/DrivingLicenseExamsDateFrames2?CategoryCode=4&CenterId=%s&ExamDate=%s", centerId, examDate)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
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

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/available-dates", availableDatesHandler)
	mux.HandleFunc("/api/available-hours", availableHoursHandler)

	port := "8080"
	fmt.Printf("Starting server on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
