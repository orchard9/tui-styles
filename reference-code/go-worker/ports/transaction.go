package ports

import (
	"context"
)

//go:generate mockgen -source=transaction.go -destination=mocks/transaction.go -package=mocks

// TransactionManager defines the contract for managing database transactions.
//
// This is a PORT (interface) in hexagonal architecture.
// Implementations are ADAPTERS (e.g., PostgreSQL transaction manager).
//
// Purpose:
// - Enable atomic operations across multiple repository calls
// - Ensure data consistency for complex business operations
// - Provide rollback on error
//
// Design Principles:
// - Transaction context propagated to all repository operations
// - Automatic rollback on error
// - Automatic commit on success
// - Nested transactions not supported (reuse outer transaction)
//
// Example Implementations:
// - adapters/repository/postgres/transaction.go
//
// Thread Safety:
// - Each transaction is isolated to its context
// - TransactionManager itself MUST be safe for concurrent use
type TransactionManager interface {
	// WithTransaction executes fn within a database transaction.
	//
	// Behavior:
	// - Begins transaction before calling fn
	// - Commits transaction if fn returns nil
	// - Rolls back transaction if fn returns error
	// - Propagates transaction via context
	//
	// Context Propagation:
	// - The transaction context (txCtx) MUST be used for all repository operations
	// - Using the original ctx will bypass the transaction
	// - Repository implementations detect and use transaction from context
	//
	// Parameters:
	// - ctx: Parent context (for timeout/cancellation)
	// - fn: Function to execute within transaction
	//
	// Returns:
	// - error: Error from fn (transaction already rolled back)
	// - error: domain.ErrInternal for transaction management failures
	//
	// Example - Transfer Money (Atomic):
	//   err := txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
	//       // All operations use txCtx for transaction participation
	//       if err := accountRepo.Withdraw(txCtx, fromAccountID, amount); err != nil {
	//           return err  // Automatic rollback
	//       }
	//       if err := accountRepo.Deposit(txCtx, toAccountID, amount); err != nil {
	//           return err  // Automatic rollback
	//       }
	//       return nil  // Automatic commit
	//   })
	//
	// Example - Create User with Profile (Atomic):
	//   err := txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
	//       // Create user
	//       if err := userRepo.Create(txCtx, user); err != nil {
	//           return err
	//       }
	//       // Create profile (same transaction)
	//       profile := entities.NewProfile(user.ID())
	//       if err := profileRepo.Create(txCtx, profile); err != nil {
	//           return err  // Rollback both user and profile
	//       }
	//       return nil  // Commit both
	//   })
	//
	// Example - With Use Case:
	//   type TransferMoneyUseCase struct {
	//       txMgr       ports.TransactionManager
	//       accountRepo ports.AccountRepository
	//   }
	//
	//   func (uc *TransferMoneyUseCase) Execute(ctx context.Context, req dto.TransferRequest) error {
	//       return uc.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
	//           // All repository calls use txCtx
	//           if err := uc.accountRepo.Withdraw(txCtx, req.FromID, req.Amount); err != nil {
	//               return err
	//           }
	//           return uc.accountRepo.Deposit(txCtx, req.ToID, req.Amount)
	//       })
	//   }
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
