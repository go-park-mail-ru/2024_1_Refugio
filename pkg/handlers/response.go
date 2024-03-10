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

// handleSuccess is a utility function to handle successful responses uniformly in the API.
func handleSuccess(w http.ResponseWriter, status int, body interface{}) {
	response := Response{
		Status: status,
		Body:   body,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
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
