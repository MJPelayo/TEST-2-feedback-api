package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
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

// ErrorResponse represents an error message returned to the client
// Consistent error format for all API errors
type ErrorResponse struct {
	Error string `json:"error"` // Human-readable error description
}

// SuccessResponse represents a success message returned to the client
type SuccessResponse struct {
	Message string   `json:"message"`        // Success message
	ID      string   `json:"id,omitempty"`   // ID of created feedback (optional)
	Data    Feedback `json:"data,omitempty"` // The feedback data (optional)
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

// ========== VALIDATION FUNCTIONS ==========

// validateName checks if the name field is valid
// Rules: not empty, at least 2 characters, max 100 characters
func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name is required")
	}
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	if len(name) > 100 {
		return fmt.Errorf("name must be less than 100 characters")
	}
	return nil
}

// validateEmail checks if the email field has valid format
// Uses regex pattern for basic email validation
func validateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("email is required")
	}
	// Basic email regex pattern - matches standard email formats
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// validateSubject checks if the subject field is valid
// Rules: not empty, max 200 characters
func validateSubject(subject string) error {
	if strings.TrimSpace(subject) == "" {
		return fmt.Errorf("subject is required")
	}
	if len(subject) > 200 {
		return fmt.Errorf("subject must be less than 200 characters")
	}
	return nil
}

// validateMessage checks if the message field is valid
// Rules: not empty, max 1000 characters
func validateMessage(message string) error {
	if strings.TrimSpace(message) == "" {
		return fmt.Errorf("message is required")
	}
	if len(message) > 1000 {
		return fmt.Errorf("message must be less than 1000 characters")
	}
	return nil
}

// validateRequest runs all validation checks on a feedback request
// Returns first validation error encountered
func validateRequest(req *FeedbackRequest) error {
	if err := validateName(req.Name); err != nil {
		return err
	}
	if err := validateEmail(req.Email); err != nil {
		return err
	}
	if err := validateSubject(req.Subject); err != nil {
		return err
	}
	if err := validateMessage(req.Message); err != nil {
		return err
	}
	return nil
}

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

// ========== HANDLERS ==========

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

	// Validate all fields
	if err := validateRequest(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
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

func main() {
	// Set up the HTTP server routes
	// Apply CORS middleware to both endpoints
	http.HandleFunc("/api/feedback", enableCORS(submitFeedback))     // POST - submit new feedback
	http.HandleFunc("/api/feedback/all", enableCORS(getAllFeedback)) // GET - retrieve all feedback
	http.HandleFunc("/health", healthCheck)                          // GET - health check (no CORS needed)

	// Start the server on port 8080
	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  POST /api/feedback       - Submit feedback form")
	fmt.Println("  GET  /api/feedback/all   - Get all feedback submissions")
	fmt.Println("  GET  /health             - Health check")
	fmt.Println("\nPress Ctrl+C to stop the server")

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
