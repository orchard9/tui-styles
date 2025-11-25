// Package services contains domain services for multi-entity operations.
//
// Domain Services:
// - Coordinate operations across multiple entities
// - Contain business logic that doesn't belong to a single entity
// - Have NO external dependencies (pure domain)
// - Are stateless (no instance variables)
//
// When to use:
// - Business operation involves multiple entities
// - Logic doesn't naturally fit in any single entity
// - Need to enforce rules across entity boundaries
//
// When NOT to use:
// - Simple single-entity operations (use entity methods)
// - Persistence operations (use repositories in application layer)
// - External service calls (use application layer)
package services

import (
	"github.com/orchard9/peach/apps/email-worker/internal/domain"
	"github.com/orchard9/peach/apps/email-worker/internal/domain/entities"
)

// UserDomainService provides domain logic for multi-entity user operations.
//
// Stateless Service:
// - No instance variables
// - All data passed as parameters
// - Pure functions
//
// Example Use Cases:
// - Validate business rules across multiple entities
// - Calculate derived values using domain knowledge
// - Enforce complex business constraints
type UserDomainService struct{}

// NewUserDomainService creates a new domain service instance.
//
// Returns: Stateless service ready for use
func NewUserDomainService() *UserDomainService {
	return &UserDomainService{}
}

// ValidateEmailUniqueness checks if email is unique within context.
//
// Business Rule: Email addresses must be unique across all users.
//
// This is a domain service because:
// - Requires knowledge of ALL users (not just one entity)
// - Enforces cross-entity constraint
// - Pure business logic (no infrastructure)
//
// Parameters:
// - email: Email to check for uniqueness
// - existingUsers: Collection of existing user entities
//
// Returns:
// - error: domain.ErrConflict if email already exists
//
// Note: In practice, this would be called by application layer
// after querying repository. This is the domain logic that determines
// what "unique" means in business terms.
//
// Example:
//
//	service := services.NewUserDomainService()
//	err := service.ValidateEmailUniqueness(newEmail, existingUsers)
//	if err != nil {
//	    // Email is not unique
//	}
func (s *UserDomainService) ValidateEmailUniqueness(email string, existingUsers []*entities.User) error {
	for _, user := range existingUsers {
		if user.Email() == email {
			return domain.ErrDuplicate("email already in use")
		}
	}
	return nil
}

// CanDeleteUser determines if a user can be safely deleted.
//
// Business Rules:
// - User must exist
// - User must not have dependent entities
// - User must not have active sessions
//
// This is domain logic because it encodes business rules about
// what makes a user "deletable" from a business perspective.
//
// Parameters:
// - user: User entity to check for deletion eligibility
// - hasActiveSessions: Whether user has active sessions
// - hasDependentData: Whether user has data that would be orphaned
//
// Returns:
// - error: domain.ErrValidation if user cannot be deleted
//
// Example:
//
//	service := services.NewUserDomainService()
//	err := service.CanDeleteUser(user, hasActiveSessions, hasDependentData)
//	if err != nil {
//	    // Cannot delete user
//	}
func (s *UserDomainService) CanDeleteUser(user *entities.User, hasActiveSessions, hasDependentData bool) error {
	if hasActiveSessions {
		return domain.ErrInvalidInput("cannot delete user with active sessions")
	}

	if hasDependentData {
		return domain.ErrInvalidInput("cannot delete user with dependent data")
	}

	return nil
}

// IsValidUserTransition checks if state transition is allowed.
//
// Business Rule: Users can only transition through valid states.
//
// Example business state machine:
// - New -> Active: Always allowed
// - Active -> Suspended: Requires reason
// - Suspended -> Active: Requires approval
// - Any -> Deleted: Requires special permission
//
// This domain service encodes the business rules about valid
// state transitions.
//
// Parameters:
// - currentStatus: Current user status
// - newStatus: Desired new status
//
// Returns:
// - bool: true if transition is valid
// - error: domain.ErrValidation if transition not allowed
//
// Example:
//
//	service := services.NewUserDomainService()
//	ok, err := service.IsValidUserTransition("active", "suspended")
func (s *UserDomainService) IsValidUserTransition(currentStatus, newStatus string) (bool, error) {
	// Define valid state transitions
	validTransitions := map[string][]string{
		"new":       {"active"},
		"active":    {"suspended", "deleted"},
		"suspended": {"active", "deleted"},
		"deleted":   {}, // Cannot transition from deleted
	}

	allowedStates, exists := validTransitions[currentStatus]
	if !exists {
		return false, domain.ErrInvalidInput("invalid current status")
	}

	for _, allowed := range allowedStates {
		if allowed == newStatus {
			return true, nil
		}
	}

	return false, domain.ErrInvalidInput("invalid status transition")
}
