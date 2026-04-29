// File: handlers.go
// Purpose: Contains all HTTP handler functions for API endpoints

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Global storage instance - accessible throughout the application
// In a real app, this might be a database connection
var storage = NewStorage()

// ========== MIDDLEWARE ==========

// enableCORS adds CORS headers to allow cross-origin requests
// This is necessary for web browsers to call the API from different domains
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers to allow any origin (for development)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// ========== API HANDLERS ==========

// submitFeedback handles POST requests to /api/feedback
// Validates input, saves to storage, returns created feedback
func submitFeedback(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body
	var req FeedbackRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Reject unexpected JSON fields

	if err := decoder.Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON: " + err.Error()})
		return
	}

	// Validate all fields using our error handling system
	if validationErr := ValidateRequest(&req); validationErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: validationErr.Error()})
		return
	}

	// Create new feedback with auto-generated fields
	feedback := Feedback{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()), // Nanosecond timestamp as unique ID
		Name:      strings.TrimSpace(req.Name),
		Email:     strings.ToLower(strings.TrimSpace(req.Email)), // Normalize email to lowercase
		Subject:   strings.TrimSpace(req.Subject),
		Message:   strings.TrimSpace(req.Message),
		CreatedAt: time.Now(),
	}

	// Save to in-memory storage
	storage.Save(feedback)

	// Return success response with created data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created status
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: "Feedback submitted successfully",
		ID:      feedback.ID,
		Data:    feedback,
	})
}

// getAllFeedback handles GET requests to /api/feedback/all
// Returns all feedback submissions as a JSON array
func getAllFeedback(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve all feedback from storage
	allFeedbacks := storage.GetAll()

	// Return JSON array of all feedback entries
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK status
	json.NewEncoder(w).Encode(allFeedbacks)
}

// healthCheck handles the /health endpoint
// Returns a simple status to confirm the server is running
func healthCheck(w http.ResponseWriter, r *http.Request) {
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	// Send 200 OK status with a simple message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

// serveHTML serves the index.html file for the web interface
func serveHTML(w http.ResponseWriter, r *http.Request) {
	// Only serve the root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Set content type and serve the HTML file
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "index.html")
}
