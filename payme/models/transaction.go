// Package models defines shared types used across the Payme Business SDK,
// including transactions, accounts, and payment receivers.
package models

// TransactionState represents the state of a transaction.
type TransactionState int

const (
	// StateCreated indicates a newly created transaction.
	StateCreated TransactionState = 1

	// StateCompleted indicates a successfully completed transaction.
	StateCompleted TransactionState = 2

	// StateCancelled indicates a cancelled transaction (before completion).
	StateCancelled TransactionState = -1

	// StateCancelledAfterComplete indicates a cancelled transaction (after completion).
	StateCancelledAfterComplete TransactionState = -2
)

// CancelReason represents the reason for cancelling a transaction.
type CancelReason int

const (
	// ReasonUserInitiated indicates the user initiated the cancellation.
	ReasonUserInitiated CancelReason = 1

	// ReasonMerchantInitiated indicates the merchant initiated the cancellation.
	ReasonMerchantInitiated CancelReason = 2

	// ReasonTransactionNotFound indicates the transaction was not found.
	ReasonTransactionNotFound CancelReason = 3

	// ReasonTimeout indicates the transaction was cancelled due to timeout.
	ReasonTimeout CancelReason = 4
)

// Transaction represents a financial transaction in the Payme system.
type Transaction struct {
	// ID is the unique 24-character identifier of the transaction in Payme Business.
	ID string `json:"id"`

	// Time is the timestamp (in milliseconds) when the transaction was created in Payme Business.
	Time int64 `json:"time"`

	// Amount is the transaction amount in tiyin (1/100 of UZS).
	Amount int64 `json:"amount"`

	// Account is the payer's account information.
	Account Account `json:"account"`

	// CreateTime is the timestamp (in milliseconds) when the transaction was created in the merchant's system.
	CreateTime int64 `json:"create_time"`

	// PerformTime is the timestamp (in milliseconds) when the transaction was performed.
	PerformTime int64 `json:"perform_time"`

	// CancelTime is the timestamp (in milliseconds) when the transaction was cancelled.
	CancelTime int64 `json:"cancel_time"`

	// Transaction is the merchant's internal transaction identifier.
	Transaction string `json:"transaction"`

	// State is the current state of the transaction.
	State TransactionState `json:"state"`

	// Reason is the cancellation reason (null if not cancelled).
	Reason *CancelReason `json:"reason"`

	// Receivers is the list of payment receivers (for chained payments).
	Receivers []Receiver `json:"receivers,omitempty"`
}
