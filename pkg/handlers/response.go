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

// enableCors adds CORS headers to the server response.
// Headers include allowed origin, methods, and headers.
// If the request is an OPTIONS preflight request, a successful status is returned without further processing.
func enableCors(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "127.0.0.1, localhost, 89.208.223.140")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusOK)
		return
	}
}
