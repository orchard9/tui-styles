// Package handlers contains HTTP handlers (adapters) for the delivery layer.
//
// Handlers are ADAPTERS in hexagonal architecture:
// - Convert HTTP requests to use case inputs
// - Execute use cases
// - Convert use case outputs to HTTP responses
// - Handle HTTP-specific concerns (status codes, headers)
//
// Responsibilities:
// - Request validation (HTTP-level)
// - Request/Response mapping
// - HTTP status code selection
// - Error response formatting
// - Route registration
//
// NOT Responsible For:
// - Business logic (belongs in domain/use cases)
// - Data persistence (belongs in repositories)
// - Complex validation (belongs in domain)
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/orchard9/peach/apps/email-worker/internal/domain"
	"github.com/orchard9/peach/apps/email-worker/internal/usecases"
)

// UserHandler handles HTTP requests for user operations
//
// Dependencies:
// - Use cases for each operation
//
// Thread-Safety: Safe for concurrent use (stateless, use cases handle concurrency)
type UserHandler struct {
	createUseCase      *usecases.CreateUserUseCase
	getUseCase         *usecases.GetUserUseCase
	updateEmailUseCase *usecases.UpdateUserEmailUseCase
	updateNameUseCase  *usecases.UpdateUserNameUseCase
	deleteUseCase      *usecases.DeleteUserUseCase
	listUseCase        *usecases.ListUsersUseCase
	searchEmailUseCase *usecases.SearchUserByEmailUseCase
}

// NewUserHandler creates a new HTTP handler
//
// Parameters:
// - All use cases for user operations
//
// Returns: Configured handler ready for route registration
func NewUserHandler(
	createUseCase *usecases.CreateUserUseCase,
	getUseCase *usecases.GetUserUseCase,
	updateEmailUseCase *usecases.UpdateUserEmailUseCase,
	updateNameUseCase *usecases.UpdateUserNameUseCase,
	deleteUseCase *usecases.DeleteUserUseCase,
	listUseCase *usecases.ListUsersUseCase,
	searchEmailUseCase *usecases.SearchUserByEmailUseCase,
) *UserHandler {
	return &UserHandler{
		createUseCase:      createUseCase,
		getUseCase:         getUseCase,
		updateEmailUseCase: updateEmailUseCase,
		updateNameUseCase:  updateNameUseCase,
		deleteUseCase:      deleteUseCase,
		listUseCase:        listUseCase,
		searchEmailUseCase: searchEmailUseCase,
	}
}

// RegisterRoutes registers all user routes
//
// Route Structure:
// - POST   /users              Create new user
// - GET    /users              List users (with pagination)
// - GET    /users/:id          Get user by ID
// - PATCH  /users/:id/email    Update user email
// - PATCH  /users/:id/name     Update user name
// - DELETE /users/:id          Delete user
//
// Query Parameters:
// - GET /users?email=<email>  Search by email
// - GET /users?limit=<n>&offset=<n>  Pagination
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Patch("/{id}/email", h.UpdateEmail)
		r.Patch("/{id}/name", h.UpdateName)
		r.Delete("/{id}", h.Delete)
	})
}

// Create handles POST /users
//
// Request Body:
//
//	{
//	  "email": "user@example.com",
//	  "name": "John Doe"
//	}
//
// Response: 201 Created
//
//	{
//	  "id": "uuid",
//	  "email": "user@example.com",
//	  "name": "John Doe"
//	}
//
// Error Responses:
// - 400 Bad Request: Invalid input, malformed JSON
// - 409 Conflict: Email already exists
// - 500 Internal Server Error: Unexpected error
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Decode request
	var req struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Execute use case
	input := usecases.CreateUserInput{
		Email: req.Email,
		Name:  req.Name,
	}
	output, err := h.createUseCase.Execute(r.Context(), input)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	// Respond with created resource
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(output)
}

// Get handles GET /users/:id
//
// Response: 200 OK
//
//	{
//	  "id": "uuid",
//	  "email": "user@example.com",
//	  "name": "John Doe",
//	  "created_at": "2025-01-01T00:00:00Z",
//	  "updated_at": "2025-01-01T00:00:00Z"
//	}
//
// Error Responses:
// - 400 Bad Request: Invalid ID format
// - 404 Not Found: User doesn't exist
// - 500 Internal Server Error: Unexpected error
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing ID parameter")
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Execute use case
	input := usecases.GetUserInput{ID: id}
	output, err := h.getUseCase.Execute(r.Context(), input)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	// Respond with resource
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(output)
}

// List handles GET /users
//
// Query Parameters:
// - limit: Max results per page (default: 20, max: 100)
// - offset: Number of results to skip (default: 0)
// - email: Search by email (optional)
//
// Response: 200 OK
//
//	{
//	  "users": [...],
//	  "total": 100,
//	  "limit": 20,
//	  "offset": 0
//	}
//
// Error Responses:
// - 400 Bad Request: Invalid query parameters
// - 500 Internal Server Error: Unexpected error
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	// Check if searching by email
	if email := r.URL.Query().Get("email"); email != "" {
		h.SearchByEmail(w, r, email)
		return
	}

	// Parse pagination parameters
	limit, offset := parsePagination(r)

	// Execute use case
	input := usecases.ListUsersInput{
		Limit:  limit,
		Offset: offset,
	}
	output, err := h.listUseCase.Execute(r.Context(), input)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	// Respond with list
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(output)
}

// SearchByEmail handles GET /users?email=<email>
func (h *UserHandler) SearchByEmail(w http.ResponseWriter, r *http.Request, email string) {
	// Execute use case
	input := usecases.SearchUserByEmailInput{Email: email}
	output, err := h.searchEmailUseCase.Execute(r.Context(), input)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	// Respond with single result (wrapped in list format for consistency)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"users":  []interface{}{output},
		"total":  1,
		"limit":  1,
		"offset": 0,
	})
}

// UpdateEmail handles PATCH /users/:id/email
//
// Request Body:
//
//	{
//	  "email": "newemail@example.com"
//	}
//
// Response: 200 OK
//
//	{
//	  "id": "uuid",
//	  "email": "newemail@example.com"
//	}
//
// Error Responses:
// - 400 Bad Request: Invalid ID or email
// - 404 Not Found: User doesn't exist
// - 409 Conflict: Email already exists
// - 500 Internal Server Error: Unexpected error
func (h *UserHandler) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing ID parameter")
		return
	}

	// Decode request
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Execute use case
	input := usecases.UpdateUserEmailInput{
		ID:       id,
		NewEmail: req.Email,
	}
	output, err := h.updateEmailUseCase.Execute(r.Context(), input)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	// Respond with updated resource
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(output)
}

// UpdateName handles PATCH /users/:id/name
//
// Request Body:
//
//	{
//	  "name": "Jane Doe"
//	}
//
// Response: 200 OK
//
//	{
//	  "id": "uuid",
//	  "name": "Jane Doe"
//	}
func (h *UserHandler) UpdateName(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing ID parameter")
		return
	}

	// Decode request
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Execute use case
	input := usecases.UpdateUserNameInput{
		ID:      id,
		NewName: req.Name,
	}
	output, err := h.updateNameUseCase.Execute(r.Context(), input)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	// Respond with updated resource
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(output)
}

// Delete handles DELETE /users/:id
//
// Response: 204 No Content
//
// Error Responses:
// - 400 Bad Request: Invalid ID format
// - 500 Internal Server Error: Unexpected error
//
// Note: Idempotent - returns 204 even if user doesn't exist
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing ID parameter")
		return
	}

	// Execute use case
	input := usecases.DeleteUserInput{ID: id}
	_, err := h.deleteUseCase.Execute(r.Context(), input)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	// Respond with no content (idempotent delete)
	w.WriteHeader(http.StatusNoContent)
}

// Helper functions

// handleUseCaseError maps domain errors to HTTP status codes
func handleUseCaseError(w http.ResponseWriter, err error) {
	var status int
	var message string

	switch {
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
		message = err.Error()
	case errors.Is(err, domain.ErrConflict):
		status = http.StatusConflict
		message = err.Error()
	case errors.Is(err, domain.ErrInvalidInput("")):
		status = http.StatusBadRequest
		message = err.Error()
	default:
		status = http.StatusInternalServerError
		message = "Internal server error"
	}

	respondError(w, status, message)
}

// respondError writes error response
func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// parsePagination extracts pagination parameters from request
func parsePagination(r *http.Request) (limit, offset int) {
	// Default values
	limit = 20
	offset = 0

	// Parse limit
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
			if limit > 100 {
				limit = 100 // Max limit
			}
		}
	}

	// Parse offset
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return limit, offset
}
