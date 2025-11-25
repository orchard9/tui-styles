// Package usecases contains application use cases.
//
// Use cases orchestrate domain logic and coordinate between ports.
// They represent specific business operations in a vertical slice architecture.
//
// Responsibilities:
// - Coordinate domain entities and services
// - Orchestrate calls to ports (repositories, cache, external services)
// - Handle transaction boundaries
// - Validate inputs and outputs
// - Map between DTOs and domain entities
//
// Rules:
// - One use case per business operation (vertical slice)
// - Use cases should be thin orchestrators, not contain business logic
// - Business logic belongs in domain entities/services
// - Use context.Context for cancellation, timeouts, and transactions
// - Return domain errors, not infrastructure errors
package usecases

import (
	"context"

	"github.com/orchard9/peach/apps/email-worker/internal/domain"
	"github.com/orchard9/peach/apps/email-worker/internal/domain/entities"
	"github.com/orchard9/peach/apps/email-worker/internal/ports"
)

// CreateUserInput contains the data needed to create a user
//
// This is a DTO (Data Transfer Object) that decouples the use case
// from HTTP handlers, gRPC services, or other delivery mechanisms.
type CreateUserInput struct {
	Email string
	Name  string
}

// CreateUserOutput contains the result of creating a user
//
// Returns only what the caller needs, not the full entity.
// This maintains encapsulation and allows us to change internals.
type CreateUserOutput struct {
	ID    string
	Email string
	Name  string
}

// CreateUserUseCase handles the business operation of creating a user
//
// Single Responsibility: Only creates users
//
// Dependencies:
// - UserRepository: For persisting the user
//
// Thread-Safety: Safe for concurrent use (stateless)
type CreateUserUseCase struct {
	repo ports.UserRepository
}

// NewCreateUserUseCase creates a new use case instance
//
// Parameters:
// - repo: Repository for user persistence
//
// Example:
//
//	useCase := usecases.NewCreateUserUseCase(postgresRepo)
//	output, err := useCase.Execute(ctx, input)
func NewCreateUserUseCase(repo ports.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		repo: repo,
	}
}

// Execute performs the create user operation
//
// Process:
// 1. Validate input
// 2. Create domain entity (enforces business rules)
// 3. Persist via repository
// 4. Return output DTO
//
// Parameters:
// - ctx: Context for cancellation, timeouts, and tracing
// - input: Data needed to create the user
//
// Returns:
// - CreateUserOutput: The created user data
// - error: Domain error if operation fails
//
// Possible Errors:
// - domain.ErrInvalidInput: Invalid email or name
// - domain.ErrConflict: Email already exists
// - domain.ErrInternal: Database or infrastructure error
//
// Example:
//
//	input := usecases.CreateUserInput{
//	    Email: "user@example.com",
//	    Name:  "John Doe",
//	}
//	output, err := useCase.Execute(ctx, input)
//	if err != nil {
//	    // Handle error
//	}
//	fmt.Printf("Created user with ID: %s\n", output.ID)
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	// Step 1: Validate input
	if input.Email == "" {
		return nil, domain.ErrInvalidInput("email is required")
	}
	if input.Name == "" {
		return nil, domain.ErrInvalidInput("name is required")
	}

	// Step 2: Create domain entity (domain enforces business rules)
	entity, err := entities.NewUser(input.Email, input.Name)
	if err != nil {
		return nil, err
	}

	// Step 3: Persist via repository
	if err := uc.repo.Create(ctx, entity); err != nil {
		return nil, err
	}

	// Step 4: Return output DTO
	return &CreateUserOutput{
		ID:    entity.ID().String(),
		Email: entity.Email(),
		Name:  entity.Name(),
	}, nil
}
