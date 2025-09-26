package api

import (
	"net/http"
)

// RegisterRoutes registers the API routes.
func (a *API) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/health", a.HealthHandler)
	mux.HandleFunc("/api/available-dates", a.AvailableDatesHandler)
	mux.HandleFunc("/api/available-hours", a.AvailableHoursHandler)
	mux.HandleFunc("/api/city-info", a.CityInfoHandler)
}
