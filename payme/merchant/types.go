package merchant

import "github.com/xusk947/payme-sdk-go/payme/models"

// CheckPerformTransactionRequest is the request for the CheckPerformTransaction method.
// It checks whether a financial transaction can be created.
type CheckPerformTransactionRequest struct {
	// ID is the unique transaction identifier in Payme Business (24-character hex string).
	ID string `json:"id"`

	// Time is the timestamp in milliseconds when the request was created.
	Time int64 `json:"time"`

	// Amount is the payment amount in tiyin.
	Amount int64 `json:"amount"`

	// Account is the payer's account information.
	Account models.Account `json:"account"`
}

// CheckPerformTransactionResponse is the response for the CheckPerformTransaction method.
type CheckPerformTransactionResponse struct {
	// Allow indicates whether the transaction is allowed.
	Allow bool `json:"allow"`

	// Additional is optional additional information returned to the Payme app.
	Additional map[string]any `json:"additional,omitempty"`
}

// CreateTransactionRequest is the request for the CreateTransaction method.
type CreateTransactionRequest struct {
	// ID is the unique transaction identifier in Payme Business.
	ID string `json:"id"`

	// Time is the timestamp in milliseconds when the transaction was created in Payme Business.
	Time int64 `json:"time"`

	// Amount is the payment amount in tiyin.
	Amount int64 `json:"amount"`

	// Account is the payer's account information.
	Account models.Account `json:"account"`
}

// CreateTransactionResponse is the response for the CreateTransaction method.
type CreateTransactionResponse struct {
	// CreateTime is the timestamp in milliseconds when the transaction was created in the merchant's system.
	CreateTime int64 `json:"create_time"`

	// Transaction is the merchant's internal transaction identifier.
	Transaction string `json:"transaction"`

	// State is the transaction state (1 = created).
	State models.TransactionState `json:"state"`

	// Receivers is the list of payment receivers for chained payments.
	// Can be nil or omitted for direct payments.
	Receivers []models.Receiver `json:"receivers,omitempty"`
}

// PerformTransactionRequest is the request for the PerformTransaction method.
type PerformTransactionRequest struct {
	// ID is the unique transaction identifier in Payme Business.
	ID string `json:"id"`
}

// PerformTransactionResponse is the response for the PerformTransaction method.
type PerformTransactionResponse struct {
	// Transaction is the merchant's internal transaction identifier.
	Transaction string `json:"transaction"`

	// PerformTime is the timestamp in milliseconds when the transaction was performed.
	PerformTime int64 `json:"perform_time"`

	// State is the transaction state (2 = completed).
	State models.TransactionState `json:"state"`
}

// CancelTransactionRequest is the request for the CancelTransaction method.
type CancelTransactionRequest struct {
	// ID is the unique transaction identifier in Payme Business.
	ID string `json:"id"`

	// Reason is the cancellation reason code.
	Reason models.CancelReason `json:"reason"`
}

// CancelTransactionResponse is the response for the CancelTransaction method.
type CancelTransactionResponse struct {
	// Transaction is the merchant's internal transaction identifier.
	Transaction string `json:"transaction"`

	// CancelTime is the timestamp in milliseconds when the transaction was cancelled.
	CancelTime int64 `json:"cancel_time"`

	// State is the transaction state (-1 = cancelled, -2 = cancelled after complete).
	State models.TransactionState `json:"state"`
}

// CheckTransactionRequest is the request for the CheckTransaction method.
type CheckTransactionRequest struct {
	// ID is the unique transaction identifier in Payme Business.
	ID string `json:"id"`
}

// CheckTransactionResponse is the response for the CheckTransaction method.
type CheckTransactionResponse struct {
	// CreateTime is the timestamp in milliseconds when the transaction was created.
	CreateTime int64 `json:"create_time"`

	// PerformTime is the timestamp in milliseconds when the transaction was performed.
	PerformTime int64 `json:"perform_time"`

	// CancelTime is the timestamp in milliseconds when the transaction was cancelled.
	CancelTime int64 `json:"cancel_time"`

	// Transaction is the merchant's internal transaction identifier.
	Transaction string `json:"transaction"`

	// State is the transaction state.
	State models.TransactionState `json:"state"`

	// Reason is the cancellation reason (nil if not cancelled).
	Reason *models.CancelReason `json:"reason"`
}

// GetStatementRequest is the request for the GetStatement method.
type GetStatementRequest struct {
	// From is the start timestamp in milliseconds of the period.
	From int64 `json:"from"`

	// To is the end timestamp in milliseconds of the period.
	To int64 `json:"to"`
}

// GetStatementResponse is the response for the GetStatement method.
type GetStatementResponse struct {
	// Transactions is the list of transactions for the specified period.
	Transactions []models.Transaction `json:"transactions"`
}
