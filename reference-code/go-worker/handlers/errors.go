// Package handlers provides HTTP error handling utilities.
package handlers

// ErrorResponse represents a standardized API error response
//
// Format:
//
//	{
//	  "error": "Human-readable error message",
//	  "code": "ERROR_CODE",
//	  "fields": {
//	    "email": "must be a valid email address",
//	    "name": "is required"
//	  }
//	}
type ErrorResponse struct {
	Error  string            `json:"error"`
	Code   string            `json:"code,omitempty"`
	Fields map[string]string `json:"fields,omitempty"`
}
