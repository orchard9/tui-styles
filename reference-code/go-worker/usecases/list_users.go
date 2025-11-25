package usecases

import (
	"context"
	"time"

	"github.com/orchard9/peach/apps/email-worker/internal/domain"
	"github.com/orchard9/peach/apps/email-worker/internal/ports"
)

// ListUsersInput contains the data needed to list users
//
// Implements pagination to handle large datasets efficiently.
type ListUsersInput struct {
	Limit  int // Max results per page (validated: 1-100)
	Offset int // Number of results to skip
}

// ListUsersOutput contains the list of users
type ListUsersOutput struct {
	Users  []UserItem
	Total  int // Total count of users (for pagination UI)
	Limit  int // Limit used for this query
	Offset int // Offset used for this query
}

// UserItem represents a user in the list
//
// Note: This is a subset of the full entity, optimized for list views.
// Clients can fetch full details using GetUserUseCase.
type UserItem struct {
	ID        string
	Email     string
	Name      string
	CreatedAt string
}

// ListUsersUseCase handles the business operation of listing users
//
// Single Responsibility: Only lists users with pagination
//
// Dependencies:
// - UserRepository: For fetching user list
//
// Thread-Safety: Safe for concurrent use (stateless)
type ListUsersUseCase struct {
	repo ports.UserRepository
}

// NewListUsersUseCase creates a new use case instance
//
// Parameters:
// - repo: Repository for user listing
func NewListUsersUseCase(repo ports.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{
		repo: repo,
	}
}

// Execute performs the list users operation
//
// Process:
// 1. Validate and normalize input
// 2. Fetch from repository
// 3. Map entities to output DTOs
// 4. Return output with metadata
//
// Parameters:
// - ctx: Context for cancellation, timeouts, and tracing
// - input: Pagination parameters
//
// Returns:
// - ListUsersOutput: List of users with pagination metadata
// - error: Domain error if operation fails
//
// Possible Errors:
// - domain.ErrInvalidInput: Invalid pagination parameters
// - domain.ErrInternal: Database or infrastructure error
//
// Pagination Rules:
// - Limit: 1-100 (default: 20)
// - Offset: >= 0 (default: 0)
// - Results ordered by created_at DESC (newest first)
//
// Example:
//
//	input := usecases.ListUsersInput{
//	    Limit:  20,
//	    Offset: 0,
//	}
//	output, err := useCase.Execute(ctx, input)
//	for _, user := range output.Users {
//	    fmt.Printf("%s: %s\n", user.Name, user.Email)
//	}
func (uc *ListUsersUseCase) Execute(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	// Step 1: Validate and normalize input
	limit := input.Limit
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		return nil, domain.ErrInvalidInput("limit cannot exceed 100")
	}

	offset := input.Offset
	if offset < 0 {
		return nil, domain.ErrInvalidInput("offset cannot be negative")
	}

	// Step 2: Fetch from repository using pagination options
	listOpts := ports.ListOptions{
		Limit:  limit,
		Offset: offset,
	}

	entities, err := uc.repo.List(ctx, listOpts)
	if err != nil {
		return nil, err
	}

	// Step 3: Map entities to output DTOs
	items := make([]UserItem, 0, len(entities))
	for _, entity := range entities {
		items = append(items, UserItem{
			ID:        entity.ID().String(),
			Email:     entity.Email(),
			Name:      entity.Name(),
			CreatedAt: entity.CreatedAt().Format(time.RFC3339),
		})
	}

	// Step 4: Return output with metadata
	// Note: Total count would require a separate COUNT query
	// For now, we return len(items) as a simple implementation
	return &ListUsersOutput{
		Users:  items,
		Total:  len(items), // TODO: Implement proper total count if needed
		Limit:  limit,
		Offset: offset,
	}, nil
}

// SearchUserByEmailInput contains the data needed to search users by email
type SearchUserByEmailInput struct {
	Email string
}

// SearchUserByEmailOutput contains the search result
type SearchUserByEmailOutput struct {
	ID        string
	Email     string
	Name      string
	CreatedAt string
	UpdatedAt string
}

// SearchUserByEmailUseCase handles searching users by email
//
// Single Responsibility: Only searches users by email address
type SearchUserByEmailUseCase struct {
	repo ports.UserRepository
}

// NewSearchUserByEmailUseCase creates a new use case instance
func NewSearchUserByEmailUseCase(repo ports.UserRepository) *SearchUserByEmailUseCase {
	return &SearchUserByEmailUseCase{
		repo: repo,
	}
}

// Execute performs the search user by email operation
func (uc *SearchUserByEmailUseCase) Execute(
	ctx context.Context,
	input SearchUserByEmailInput,
) (*SearchUserByEmailOutput, error) {
	// Step 1: Validate input
	if input.Email == "" {
		return nil, domain.ErrInvalidInput("email is required")
	}

	// Step 2: Search via repository
	entity, err := uc.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	// Step 3: Return output DTO
	return &SearchUserByEmailOutput{
		ID:        entity.ID().String(),
		Email:     entity.Email(),
		Name:      entity.Name(),
		CreatedAt: entity.CreatedAt().Format(time.RFC3339),
		UpdatedAt: entity.UpdatedAt().Format(time.RFC3339),
	}, nil
}
