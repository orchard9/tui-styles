package usecases

import (
	"context"

	"github.com/google/uuid"

	"github.com/orchard9/peach/apps/email-worker/internal/domain"
	"github.com/orchard9/peach/apps/email-worker/internal/ports"
)

// DeleteUserInput contains the data needed to delete a user
type DeleteUserInput struct {
	ID string
}

// DeleteUserOutput contains the result of deleting a user
//
// Note: We return the ID to confirm which entity was deleted.
// This is useful for audit logging and client-side state management.
type DeleteUserOutput struct {
	ID string
}

// DeleteUserUseCase handles the business operation of deleting a user
//
// Single Responsibility: Only deletes users
//
// Dependencies:
// - UserRepository: For deleting the user
//
// Thread-Safety: Safe for concurrent use (stateless)
type DeleteUserUseCase struct {
	repo ports.UserRepository
}

// NewDeleteUserUseCase creates a new use case instance
//
// Parameters:
// - repo: Repository for user deletion
func NewDeleteUserUseCase(repo ports.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		repo: repo,
	}
}

// Execute performs the delete user operation
//
// Process:
// 1. Validate input
// 2. Delete from repository (idempotent)
// 3. Invalidate cache
// 4. Return output DTO
//
// Parameters:
// - ctx: Context for cancellation, timeouts, and transactions
// - input: ID of the user to delete
//
// Returns:
// - DeleteUserOutput: Confirmation of deletion
// - error: Domain error if operation fails
//
// Possible Errors:
// - domain.ErrInvalidInput: Invalid ID format
// - domain.ErrInternal: Database or infrastructure error
//
// Note: This operation is idempotent. Deleting a non-existent user
// is not an error (follows REST DELETE semantics).
//
// Example:
//
//	input := usecases.DeleteUserInput{
//	    ID: "550e8400-e29b-41d4-a716-446655440000",
//	}
//	output, err := useCase.Execute(ctx, input)
func (uc *DeleteUserUseCase) Execute(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, error) {
	// Step 1: Validate input
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, domain.ErrInvalidInput("invalid user ID format")
	}

	// Step 2: Delete from repository (idempotent operation)
	if err := uc.repo.Delete(ctx, id); err != nil {
		return nil, err
	}

	// Step 3: Return output DTO
	return &DeleteUserOutput{
		ID: id.String(),
	}, nil
}
