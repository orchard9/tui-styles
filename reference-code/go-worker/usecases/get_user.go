package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/orchard9/peach/apps/email-worker/internal/domain"
	"github.com/orchard9/peach/apps/email-worker/internal/domain/entities"
	"github.com/orchard9/peach/apps/email-worker/internal/ports"
)

// GetUserInput contains the data needed to retrieve a user
type GetUserInput struct {
	ID string
}

// GetUserOutput contains the retrieved user data
type GetUserOutput struct {
	ID        string
	Email     string
	Name      string
	CreatedAt string
	UpdatedAt string
}

// GetUserUseCase handles the business operation of retrieving a user
//
// Single Responsibility: Only retrieves users by ID with caching
//
// Dependencies:
// - UserRepository: For fetching the user
// - Cache: Optional cache for performance optimization
//
// Caching Strategy:
// - Cache-aside pattern: Check cache first, fetch from DB on miss
// - Fail open: Cache failures don't prevent data retrieval
// - Specific cache keys: user:id:<uuid>
// - TTL: 5 minutes
//
// Thread-Safety: Safe for concurrent use (stateless)
type GetUserUseCase struct {
	repo  ports.UserRepository
	cache ports.CacheService
}

// NewGetUserUseCase creates a new use case instance
//
// Parameters:
// - repo: Repository for user retrieval
// - cache: Optional cache (can be nil)
//
// Example:
//
//	useCase := usecases.NewGetUserUseCase(repo, cache)
//	useCase := usecases.NewGetUserUseCase(repo, nil) // Without cache
func NewGetUserUseCase(repo ports.UserRepository, cache ports.CacheService) *GetUserUseCase {
	return &GetUserUseCase{
		repo:  repo,
		cache: cache,
	}
}

// Execute performs the get user operation with caching
//
// Process:
// 1. Validate input
// 2. Check cache (if available)
// 3. Fetch from repository on cache miss
// 4. Store in cache (if available)
// 5. Return output DTO
//
// Caching Implementation:
// - Cache key: user:id:<uuid>
// - TTL: 5 minutes
// - Fail open: Cache failures logged but don't block operation
// - Serialization: JSON
//
// Parameters:
// - ctx: Context for cancellation, timeouts, and tracing
// - input: ID of the user to retrieve
//
// Returns:
// - GetUserOutput: The retrieved user data
// - error: Domain error if operation fails
//
// Possible Errors:
// - domain.ErrInvalidInput: Invalid ID format
// - domain.ErrNotFound: User doesn't exist
// - domain.ErrInternal: Database or infrastructure error
func (uc *GetUserUseCase) Execute(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	// Step 1: Validate input
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, domain.ErrInvalidInput("invalid user ID format")
	}

	// Step 2: Try cache first (cache-aside pattern)
	if uc.cache != nil {
		if cached, err := uc.getFromCache(ctx, id); err == nil && cached != nil {
			return cached, nil
		}
		// Cache miss or error - continue to repository (fail open)
	}

	// Step 3: Fetch from repository
	entity, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Step 4: Convert to output DTO
	output := uc.entityToOutput(entity)

	// Step 5: Store in cache (fail open - don't block on cache errors)
	if uc.cache != nil {
		_ = uc.setInCache(ctx, id, output)
	}

	return output, nil
}

// entityToOutput converts a domain entity to output DTO
func (uc *GetUserUseCase) entityToOutput(entity *entities.User) *GetUserOutput {
	return &GetUserOutput{
		ID:        entity.ID().String(),
		Email:     entity.Email(),
		Name:      entity.Name(),
		CreatedAt: entity.CreatedAt().Format(time.RFC3339),
		UpdatedAt: entity.UpdatedAt().Format(time.RFC3339),
	}
}

// cacheKey generates a specific cache key for a user ID
//
// Format: user:id:<uuid>
// Example: user:id:550e8400-e29b-41d4-a716-446655440000
//
// Why not wildcards:
// - Specific keys allow targeted invalidation
// - Prevents cache pollution
// - Makes debugging easier
func (uc *GetUserUseCase) cacheKey(id uuid.UUID) string {
	return fmt.Sprintf("user:id:%s", id.String())
}

// getFromCache attempts to retrieve a user from cache
//
// Returns:
// - *GetUserOutput: The cached data if found
// - error: Cache miss or deserialization error
//
// Note: Callers should treat errors as cache misses (fail open)
func (uc *GetUserUseCase) getFromCache(ctx context.Context, id uuid.UUID) (*GetUserOutput, error) {
	key := uc.cacheKey(id)

	data, err := uc.cache.Get(ctx, key)
	if err != nil {
		return nil, err // Cache miss or error
	}

	var output GetUserOutput
	if err := json.Unmarshal([]byte(data), &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached user: %w", err)
	}

	return &output, nil
}

// setInCache stores a user in cache with 5 minute TTL
//
// TTL Rationale:
// - 5 minutes balances freshness with cache hit rate
// - Short enough to prevent stale data issues
// - Long enough to absorb read spikes
//
// Returns:
// - error: Serialization or cache storage error
//
// Note: Callers should ignore errors (fail open)
func (uc *GetUserUseCase) setInCache(ctx context.Context, id uuid.UUID, output *GetUserOutput) error {
	key := uc.cacheKey(id)

	data, err := json.Marshal(output)
	if err != nil {
		return fmt.Errorf("failed to marshal user for cache: %w", err)
	}

	ttl := 5 * time.Minute
	return uc.cache.Set(ctx, key, string(data), ttl)
}
