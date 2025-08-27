package main

import (
	"fmt"
	"log"
	"net/http"

	"driving-license-city-api/internal/api"
	"driving-license-city-api/internal/service"
	"driving-license-city-api/pkg/httpclient"
)

func main() {
	// Create a new HTTP client.
	client := httpclient.New()

	// Create a new SA service.
	saService := service.NewSAService(client)

	// Create a new API.
	api := api.NewAPI(saService)

	// Create a new ServeMux.
	mux := http.NewServeMux()

	// Register the API routes.
	api.RegisterRoutes(mux)

	// Start the server.
	port := "8080"
	fmt.Printf("Starting server on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
