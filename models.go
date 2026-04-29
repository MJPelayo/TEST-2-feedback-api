// File: models.go
// Purpose: Contains all data structures and error handling types

package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ========== DATA STRUCTURES ==========

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

// ========== ERROR HANDLING TYPES ==========

// ValidationError represents a validation failure
type ValidationError struct {
	Field   string // Which field failed validation
	Message string // Human-readable error message
}

// Error implements the error interface for ValidationError
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// ========== STORAGE METHODS ==========

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

// ========== VALIDATION FUNCTIONS (ERROR HANDLING) ==========

// validateName checks if the name field is valid
// Rules: not empty, at least 2 characters, max 100 characters
func validateName(name string) *ValidationError {
	if strings.TrimSpace(name) == "" {
		return NewValidationError("name", "name is required")
	}
	if len(name) < 2 {
		return NewValidationError("name", "name must be at least 2 characters")
	}
	if len(name) > 100 {
		return NewValidationError("name", "name must be less than 100 characters")
	}
	return nil
}

// validateEmail checks if the email field has valid format
// Uses regex pattern for basic email validation
func validateEmail(email string) *ValidationError {
	if strings.TrimSpace(email) == "" {
		return NewValidationError("email", "email is required")
	}
	// Basic email regex pattern - matches standard email formats
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return NewValidationError("email", "invalid email format")
	}
	return nil
}

// validateSubject checks if the subject field is valid
// Rules: not empty, max 200 characters
func validateSubject(subject string) *ValidationError {
	if strings.TrimSpace(subject) == "" {
		return NewValidationError("subject", "subject is required")
	}
	if len(subject) > 200 {
		return NewValidationError("subject", "subject must be less than 200 characters")
	}
	return nil
}

// validateMessage checks if the message field is valid
// Rules: not empty, max 1000 characters
func validateMessage(message string) *ValidationError {
	if strings.TrimSpace(message) == "" {
		return NewValidationError("message", "message is required")
	}
	if len(message) > 1000 {
		return NewValidationError("message", "message must be less than 1000 characters")
	}
	return nil
}

// ValidateRequest runs all validation checks on a feedback request
// Returns first validation error encountered, or nil if all valid
func ValidateRequest(req *FeedbackRequest) *ValidationError {
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
