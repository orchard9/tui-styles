package usecases

import (
	"context"

	"github.com/google/uuid"

	"github.com/orchard9/peach/apps/email-worker/internal/domain"
	"github.com/orchard9/peach/apps/email-worker/internal/ports"
)

// UpdateUserEmailInput contains the data needed to update a user's email
//
// Note: This is a specific use case for email updates.
// Other field updates would have their own use cases (Single Responsibility Principle).
type UpdateUserEmailInput struct {
	ID       string
	NewEmail string
}

// UpdateUserEmailOutput contains the result of updating a user's email
type UpdateUserEmailOutput struct {
	ID    string
	Email string
}

// UpdateUserEmailUseCase handles the business operation of updating a user's email
//
// Single Responsibility: Only updates user email addresses
//
// Why separate use case?
// - Different validation rules
// - Different authorization requirements
// - Different side effects (email verification, notifications)
// - Easier to test and maintain
//
// Dependencies:
// - UserRepository: For fetching and updating the user
// - CacheService (optional): For cache invalidation after updates
//
// Thread-Safety: Safe for concurrent use (stateless)
type UpdateUserEmailUseCase struct {
	repo  ports.UserRepository
	cache ports.CacheService // Optional: nil-safe
}

// NewUpdateUserEmailUseCase creates a new use case instance
//
// Parameters:
// - repo: Repository for user updates
// - cache: Optional cache for invalidation (can be nil)
func NewUpdateUserEmailUseCase(repo ports.UserRepository, cache ports.CacheService) *UpdateUserEmailUseCase {
	return &UpdateUserEmailUseCase{
		repo:  repo,
		cache: cache,
	}
}

// Execute performs the update user email operation
//
// Process:
// 1. Validate input
// 2. Fetch existing user
// 3. Update email (domain validates)
// 4. Persist changes
// 5. Invalidate cache
// 6. Return output DTO
//
// Parameters:
// - ctx: Context for cancellation, timeouts, and transactions
// - input: Data needed to update the email
//
// Returns:
// - UpdateUserEmailOutput: The updated user data
// - error: Domain error if operation fails
//
// Possible Errors:
// - domain.ErrInvalidInput: Invalid ID or email format
// - domain.ErrNotFound: User doesn't exist
// - domain.ErrConflict: New email already exists
// - domain.ErrInternal: Database or infrastructure error
//
// Example:
//
//	input := usecases.UpdateUserEmailInput{
//	    ID:       "550e8400-e29b-41d4-a716-446655440000",
//	    NewEmail: "newemail@example.com",
//	}
//	output, err := useCase.Execute(ctx, input)
func (uc *UpdateUserEmailUseCase) Execute(
	ctx context.Context,
	input UpdateUserEmailInput,
) (*UpdateUserEmailOutput, error) {
	// Step 1: Validate input
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, domain.ErrInvalidInput("invalid user ID format")
	}

	if input.NewEmail == "" {
		return nil, domain.ErrInvalidInput("new email is required")
	}

	// Step 2: Fetch existing user
	entity, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store old email for cache invalidation
	oldEmail := entity.Email()

	// Step 3: Update email (domain entity validates the new email)
	if err := entity.ChangeEmail(input.NewEmail); err != nil {
		return nil, err
	}

	// Step 4: Persist changes
	if err := uc.repo.Update(ctx, entity); err != nil {
		return nil, err
	}

	// Step 5: Invalidate affected cache entries (fail-open: errors logged but not returned)
	uc.invalidateEmailUpdateCache(ctx, id, oldEmail, input.NewEmail)

	// Step 6: Return output DTO
	return &UpdateUserEmailOutput{
		ID:    entity.ID().String(),
		Email: entity.Email(),
	}, nil
}

// invalidateEmailUpdateCache invalidates specific cache entries affected by email update
//
// Strategy: Surgical invalidation to minimize cache thrashing
// - Invalidate by ID (primary lookup)
// - Invalidate by old email (if email-based lookups exist)
// - Invalidate by new email (if email-based lookups exist)
// - Invalidate first few pages of list results (limited blast radius)
//
// Fail-Open: Cache errors are logged but don't fail the operation
// Rationale: Cache is performance optimization, not correctness requirement
func (uc *UpdateUserEmailUseCase) invalidateEmailUpdateCache(
	ctx context.Context,
	id uuid.UUID,
	oldEmail, newEmail string,
) {
	if uc.cache == nil {
		return // No cache configured
	}

	// Build list of specific keys to invalidate
	keysToInvalidate := []string{
		"user:id:" + id.String(),    // By ID lookup
		"user:email:" + oldEmail,    // Old email lookup
		"user:email:" + newEmail,    // New email lookup (in case of recreation)
		"users:list:page:1:size:10", // First page, common size
		"users:list:page:1:size:20", // First page, alternate size
		"users:list:page:1:size:50", // First page, larger size
		"users:list:page:2:size:10", // Second page (often visible)
	}

	// Invalidate each key individually (fail-open: continue on errors)
	for _, key := range keysToInvalidate {
		if err := uc.cache.Delete(ctx, key); err != nil {
			// Log error but continue (cache invalidation should not break the operation)
			// In production, you'd use a proper logger here:
			// logger.Warn("cache invalidation failed", "key", key, "error", err)
			_ = err // Explicitly ignore
		}
	}

	// Note: We do NOT use DeletePattern here because:
	// 1. Pattern-based deletion is expensive and can cause cache storms
	// 2. Stale cache entries have TTL as fallback
	// 3. Specific key deletion provides predictable performance
	// 4. Missing a few edge-case pages is acceptable (eventual consistency)
}

// UpdateUserNameInput contains the data needed to update a user's name
type UpdateUserNameInput struct {
	ID      string
	NewName string
}

// UpdateUserNameOutput contains the result of updating a user's name
type UpdateUserNameOutput struct {
	ID   string
	Name string
}

// UpdateUserNameUseCase handles the business operation of updating a user's name
//
// Single Responsibility: Only updates user names
//
// Dependencies:
// - UserRepository: For fetching and updating the user
// - CacheService (optional): For cache invalidation after updates
type UpdateUserNameUseCase struct {
	repo  ports.UserRepository
	cache ports.CacheService // Optional: nil-safe
}

// NewUpdateUserNameUseCase creates a new use case instance
//
// Parameters:
// - repo: Repository for user updates
// - cache: Optional cache for invalidation (can be nil)
func NewUpdateUserNameUseCase(repo ports.UserRepository, cache ports.CacheService) *UpdateUserNameUseCase {
	return &UpdateUserNameUseCase{
		repo:  repo,
		cache: cache,
	}
}

// Execute performs the update user name operation
func (uc *UpdateUserNameUseCase) Execute(
	ctx context.Context,
	input UpdateUserNameInput,
) (*UpdateUserNameOutput, error) {
	// Step 1: Validate input
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, domain.ErrInvalidInput("invalid user ID format")
	}

	if input.NewName == "" {
		return nil, domain.ErrInvalidInput("new name is required")
	}

	// Step 2: Fetch existing user
	entity, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Step 3: Update name (domain entity validates)
	if err := entity.ChangeName(input.NewName); err != nil {
		return nil, err
	}

	// Step 4: Persist changes
	if err := uc.repo.Update(ctx, entity); err != nil {
		return nil, err
	}

	// Step 5: Invalidate affected cache entries (fail-open: errors logged but not returned)
	uc.invalidateNameUpdateCache(ctx, id)

	// Step 6: Return output DTO
	return &UpdateUserNameOutput{
		ID:   entity.ID().String(),
		Name: entity.Name(),
	}, nil
}

// invalidateNameUpdateCache invalidates specific cache entries affected by name update
//
// Strategy: Surgical invalidation with limited scope
// - Invalidate by ID (primary lookup)
// - Invalidate first few pages of list results (names often appear in lists)
//
// Fail-Open: Cache errors are logged but don't fail the operation
func (uc *UpdateUserNameUseCase) invalidateNameUpdateCache(ctx context.Context, id uuid.UUID) {
	if uc.cache == nil {
		return // No cache configured
	}

	// Build list of specific keys to invalidate
	keysToInvalidate := []string{
		"user:id:" + id.String(),    // By ID lookup
		"users:list:page:1:size:10", // First page, common size
		"users:list:page:1:size:20", // First page, alternate size
		"users:list:page:1:size:50", // First page, larger size
		"users:list:page:2:size:10", // Second page (often visible)
	}

	// Invalidate each key individually (fail-open: continue on errors)
	for _, key := range keysToInvalidate {
		if err := uc.cache.Delete(ctx, key); err != nil {
			// Log error but continue (cache invalidation should not break the operation)
			// In production, you'd use a proper logger here:
			// logger.Warn("cache invalidation failed", "key", key, "error", err)
			_ = err // Explicitly ignore
		}
	}
}
