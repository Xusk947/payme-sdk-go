// Package merchant implements the Payme Business Merchant API server-side
// handler, dispatching JSON-RPC 2.0 requests to a user-provided MerchantHandler.
package merchant

import "context"

// MerchantHandler is the interface that merchants must implement to handle
// Payme Business Merchant API requests. Each method corresponds to a JSON-RPC
// method called by Payme Business.
//
// Implementations must store transactions in a persistent storage and handle
// idempotency: if a transaction with the same Payme ID already exists, the
// handler should return the existing transaction's state rather than creating
// a duplicate.
type MerchantHandler interface {
	// CheckPerformTransaction checks whether a financial transaction can be created.
	// It should validate the account and amount. Return an error if the transaction
	// cannot be performed.
	CheckPerformTransaction(ctx context.Context, req *CheckPerformTransactionRequest) (*CheckPerformTransactionResponse, error)

	// CreateTransaction creates a financial transaction. The handler should:
	// - Store the transaction in persistent storage
	// - Validate the account exists
	// - Validate the amount matches the billed amount
	// - Set the order status to "awaiting payment"
	// - Return the existing transaction if it was already created (idempotency)
	CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error)

	// PerformTransaction completes a previously created transaction. The handler should:
	// - Credit the merchant's account
	// - Set the order status to "paid"
	// - Return the existing state if already performed (idempotency)
	PerformTransaction(ctx context.Context, req *PerformTransactionRequest) (*PerformTransactionResponse, error)

	// CancelTransaction cancels a transaction (either created or completed).
	// For completed transactions, this performs a refund. The handler should:
	// - Cancel the transaction and update the order status
	// - Return the existing state if already cancelled (idempotency)
	CancelTransaction(ctx context.Context, req *CancelTransactionRequest) (*CancelTransactionResponse, error)

	// CheckTransaction returns the current state of a transaction.
	CheckTransaction(ctx context.Context, req *CheckTransactionRequest) (*CheckTransactionResponse, error)

	// GetStatement returns all transactions for the specified time period.
	// Transactions should be sorted by creation time in ascending order.
	// Only transactions that were successfully created via CreateTransaction
	// should be included.
	GetStatement(ctx context.Context, req *GetStatementRequest) (*GetStatementResponse, error)
}
