package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// availableDate represents a single object in the external API's JSON array.
type availableDate struct {
	BookingDate       string `json:"bookingDate"`
	BookingDateStatus int    `json:"bookingDateStatus"`
}

// centerResult holds the outcome of an API call for a single center.
type centerResult struct {
	CenterName string          `json:"centerName"`
	Dates      []availableDate `json:"dates"`
	Error      string          `json:"error,omitempty"` // Include error information
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

// fetchDatesForCenter fetches dates for a single center and sends the result to a channel.
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

// availableDatesHandler is the HTTP handler for our REST endpoint.
func availableDatesHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for JSON response and CORS (allowing any origin for this example).
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
	// Create a new ServeMux to handle routes.
	mux := http.NewServeMux()

	// Serve static files (HTML, CSS, JS) from the "static" directory.
	// We prefix the route with /static/ and strip it for the file server.
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle the API endpoint.
	mux.HandleFunc("/api/available-dates", availableDatesHandler)

	// Handle the root route to serve index.html.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // If the path is not the root, return a 404 to avoid serving index.html for every unknown path.
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
		http.ServeFile(w, r, "./static/index.html")
	})

	port := "8080"
	fmt.Printf("Starting server on http://localhost:%s\n", port)

	// Start the HTTP server with the new ServeMux.
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
