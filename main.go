package main

import (
	"fmt"
	"log"
	"net/http"
)

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
