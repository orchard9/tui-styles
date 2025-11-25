package dto

import (
	"fmt"

	"github.com/orchard9/go-core-http-toolkit/validation"
	"github.com/orchard9/peach/apps/email-worker/internal/domain/entities"
)

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=255"`
	Email string `json:"email" validate:"required,email"`
}

// Validate validates the CreateUserRequest
func (r *CreateUserRequest) Validate() error {
	validator := validation.NewValidator()
	result := validator.Validate(r)
	if !result.Valid {
		return fmt.Errorf("validation failed: %v", result.Errors)
	}
	return nil
}

// ToDomain converts CreateUserRequest to domain entity
func (r *CreateUserRequest) ToDomain() (*entities.User, error) {
	return entities.NewUser(r.Email, r.Name)
}

// UpdateUserRequest represents the request to update an existing user
type UpdateUserRequest struct {
	Name  *string `json:"name" validate:"omitempty,min=1,max=255"`
	Email *string `json:"email" validate:"omitempty,email"`
}

// Validate validates the UpdateUserRequest
func (r *UpdateUserRequest) Validate() error {
	validator := validation.NewValidator()
	result := validator.Validate(r)
	if !result.Valid {
		return fmt.Errorf("validation failed: %v", result.Errors)
	}
	return nil
}

// ApplyTo applies the update request to a domain entity
func (r *UpdateUserRequest) ApplyTo(entity *entities.User) error {
	if r.Name != nil {
		if err := entity.ChangeName(*r.Name); err != nil {
			return err
		}
	}
	if r.Email != nil {
		if err := entity.ChangeEmail(*r.Email); err != nil {
			return err
		}
	}
	return nil
}

// UserResponse represents the response for a user
type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// FromDomain converts a domain entity to UserResponse
func (r *UserResponse) FromDomain(entity *entities.User) {
	r.ID = entity.ID().String()
	r.Name = entity.Name()
	r.Email = entity.Email()
	r.CreatedAt = entity.CreatedAt().Format("2006-01-02T15:04:05Z07:00")
	r.UpdatedAt = entity.UpdatedAt().Format("2006-01-02T15:04:05Z07:00")
}

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	Items      []*UserResponse `json:"items"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// NewCreateUserResponse creates a response from a domain entity
func NewCreateUserResponse(entity *entities.User) *UserResponse {
	resp := &UserResponse{}
	resp.FromDomain(entity)
	return resp
}

// NewUpdateUserResponse creates a response from a domain entity
func NewUpdateUserResponse(entity *entities.User) *UserResponse {
	resp := &UserResponse{}
	resp.FromDomain(entity)
	return resp
}

// NewUserResponse creates a response from a domain entity
func NewUserResponse(entity *entities.User) *UserResponse {
	resp := &UserResponse{}
	resp.FromDomain(entity)
	return resp
}

// NewUserListResponse creates a paginated list response from domain entities
func NewUserListResponse(entities []*entities.User, total int64, page, pageSize int) *UserListResponse {
	items := make([]*UserResponse, len(entities))
	for i, entity := range entities {
		items[i] = NewUserResponse(entity)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &UserListResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// ListUserRequest represents the request to list users with pagination
type ListUserRequest struct {
	Page     int    `json:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size" validate:"omitempty,min=1,max=100"`
	SortBy   string `json:"sort_by" validate:"omitempty,oneof=name created_at updated_at"`
	SortDir  string `json:"sort_dir" validate:"omitempty,oneof=asc desc"`
}

// Validate validates the ListUserRequest
func (r *ListUserRequest) Validate() error {
	validator := validation.NewValidator()
	result := validator.Validate(r)
	if !result.Valid {
		return fmt.Errorf("validation failed: %v", result.Errors)
	}
	return nil
}

// SetDefaults sets default values for pagination
func (r *ListUserRequest) SetDefaults() {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.PageSize < 1 {
		r.PageSize = 10
	}
	if r.SortBy == "" {
		r.SortBy = "created_at"
	}
	if r.SortDir == "" {
		r.SortDir = "desc"
	}
}
