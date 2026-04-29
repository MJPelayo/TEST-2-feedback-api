// File: main.go
// Purpose: Application entry point - sets up and starts the server

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Set up the HTTP server routes
	// Apply CORS middleware to API endpoints
	http.HandleFunc("/api/feedback", enableCORS(submitFeedback))     // POST - submit new feedback
	http.HandleFunc("/api/feedback/all", enableCORS(getAllFeedback)) // GET - retrieve all feedback

	// Serve the HTML form interface
	http.HandleFunc("/", serveHTML) // GET - web interface

	// Health check endpoint (no CORS needed for internal monitoring)
	http.HandleFunc("/health", healthCheck)

	// Start the server on port 8080
	port := ":8080"
	fmt.Printf("========================================\n")
	fmt.Printf("✅ Feedback API Server Started\n")
	fmt.Printf("========================================\n")
	fmt.Printf("📍 Web Interface: http://localhost%s\n", port)
	fmt.Printf("📍 API Endpoints:\n")
	fmt.Printf("   POST %s/api/feedback     - Submit feedback\n", port)
	fmt.Printf("   GET  %s/api/feedback/all - Get all feedback\n", port)
	fmt.Printf("   GET  %s/health           - Health check\n", port)
	fmt.Printf("========================================\n")
	fmt.Printf("Press Ctrl+C to stop the server\n\n")

	// Listen and serve - this blocks forever
	log.Fatal(http.ListenAndServe(port, nil))
}
