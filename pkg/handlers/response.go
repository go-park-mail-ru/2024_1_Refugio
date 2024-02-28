package handlers

import (
	"encoding/json"
	"net/http"
)

// Response represents the response format.
type Response struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

// ErrorResponse represents the error response format.
type ErrorResponse struct {
	Error string `json:"error"`
}

// handleError is a utility function to handle errors uniformly in the API responses.
func handleError(w http.ResponseWriter, status int, message string) {
	response := Response{
		Status: status,
		Body:   ErrorResponse{Error: message},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
