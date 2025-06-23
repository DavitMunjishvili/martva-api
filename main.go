package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sync"
)

var staticFiles embed.FS

// availableDate represents a single object in the external API's JSON array.
type availableDate struct {
	BookingDate       string `json:"bookingDate"`
	BookingDateStatus int    `json:"bookingDateStatus"`
}

// centerResult holds the outcome of an API call for a single center.
type centerResult struct {
	CenterName string          `json:"centerName"`
	Dates      []availableDate `json:"dates"`
	Error      string          `json:"error,omitempty"`
}

// A map of Center IDs to their corresponding city names.
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

// fetchDatesForCenter remains the same...
func fetchDatesForCenter(centerID int, centerName string, ch chan<- centerResult, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("https://api-my.sa.gov.ge/api/v1/DrivingLicensePracticalExams2/DrivingLicenseExamsDates2?CategoryCode=4&CenterId=%d", centerID)
	result := centerResult{CenterName: centerName}
	resp, err := http.Get(url)
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
	ch <- result
}

// availableDatesHandler remains the same...
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

func main() {
	mux := http.NewServeMux()

	staticContent, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("could not create sub-filesystem: %v", err)
	}

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticContent))))

	mux.HandleFunc("/api/available-dates", availableDatesHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		indexHTML, err := staticContent.ReadFile("index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})

	port := "8080"
	fmt.Printf("Starting server on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
