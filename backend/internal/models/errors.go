package models

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Errors map[string][]string `json:"errors"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Errors: map[string][]string{
			"body": {message},
		},
	}
}

// NewValidationErrorResponse creates an error response from validation errors
func NewValidationErrorResponse(validationErrors ValidationErrors) ErrorResponse {
	errorMap := make(map[string][]string)
	
	for _, ve := range validationErrors {
		errorMap[ve.Field] = append(errorMap[ve.Field], ve.Message)
	}
	
	return ErrorResponse{Errors: errorMap}
}

// WriteErrorResponse writes an error response to the HTTP response writer
func WriteErrorResponse(w http.ResponseWriter, status int, err interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	
	var response ErrorResponse
	
	switch e := err.(type) {
	case ValidationErrors:
		response = NewValidationErrorResponse(e)
	case string:
		response = NewErrorResponse(e)
	case error:
		response = NewErrorResponse(e.Error())
	default:
		response = NewErrorResponse("Internal server error")
	}
	
	json.NewEncoder(w).Encode(response)
}

// WriteJSONResponse writes a JSON response to the HTTP response writer
func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}