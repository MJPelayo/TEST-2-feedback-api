package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Feedback represents a single feedback submission from the form
// This struct maps to JSON for API responses and requests
type Feedback struct {
	ID        string    `json:"id"`         // Unique identifier for each feedback
	Name      string    `json:"name"`       // Name of the person giving feedback
	Email     string    `json:"email"`      // Email address for follow-up
	Subject   string    `json:"subject"`    // Subject/title of the feedback
	Message   string    `json:"message"`    // The actual feedback message
	CreatedAt time.Time `json:"created_at"` // Timestamp when feedback was submitted
}

// FeedbackRequest represents the expected JSON format when submitting feedback
// This is separate from Feedback because the ID and CreatedAt are auto-generated
type FeedbackRequest struct {
	Name    string `json:"name"`    // Required: person's name
	Email   string `json:"email"`   // Required: valid email address
	Subject string `json:"subject"` // Required: subject line
	Message string `json:"message"` // Required: feedback content
}

// Storage holds all feedback entries in memory
// Using a map for O(1) lookups by ID
type Storage struct {
	feedbacks map[string]Feedback // Key: feedback ID, Value: Feedback object
}

// NewStorage creates and initializes a new Storage instance
// Returns a pointer to the new Storage with an empty map
func NewStorage() *Storage {
	return &Storage{
		feedbacks: make(map[string]Feedback),
	}
}

// Save stores a feedback entry in the in-memory map
// This simulates a database for our API
func (s *Storage) Save(feedback Feedback) {
	s.feedbacks[feedback.ID] = feedback
}

// GetAll returns all feedback entries as a slice
// Converts the map values into a slice for JSON serialization
func (s *Storage) GetAll() []Feedback {
	all := make([]Feedback, 0, len(s.feedbacks))
	for _, f := range s.feedbacks {
		all = append(all, f)
	}
	return all
}

// Global storage instance - accessible throughout the application
// In a real app, this might be a database connection
var storage = NewStorage()

func main() {
	// Set up the HTTP server routes
	http.HandleFunc("/health", healthCheck)

	// Start the server on port 8080
	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Println("Press Ctrl+C to stop the server")

	// Listen and serve - this blocks forever
	log.Fatal(http.ListenAndServe(port, nil))
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
