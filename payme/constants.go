package payme

const (
	// BaseURLProd is the production endpoint for the Subscribe API.
	BaseURLProd = "https://checkout.paycom.uz/api"

	// BaseURLTest is the test/sandbox endpoint for the Subscribe API.
	BaseURLTest = "https://checkout.test.paycom.uz/api"

	// CheckoutURLProd is the production checkout URL for payment initialization.
	CheckoutURLProd = "https://paycom.uz"

	// CheckoutURLTest is the test/sandbox checkout URL for payment initialization.
	CheckoutURLTest = "https://test.paycom.uz"
)

// TransactionState represents the state of a transaction in the Merchant API.
type TransactionState int

const (
	// TransactionStateCreated indicates a newly created transaction (pending payment).
	TransactionStateCreated TransactionState = 1

	// TransactionStateCompleted indicates a successfully completed transaction.
	TransactionStateCompleted TransactionState = 2

	// TransactionStateCancelled indicates a cancelled transaction (before completion).
	TransactionStateCancelled TransactionState = -1

	// TransactionStateCancelledAfterComplete indicates a cancelled transaction (after completion, i.e. refund).
	TransactionStateCancelledAfterComplete TransactionState = -2
)

// CancelReason represents the reason for cancelling a transaction.
type CancelReason int

const (
	// CancelReasonUserInitiated indicates the user initiated the cancellation.
	CancelReasonUserInitiated CancelReason = 1

	// CancelReasonMerchantInitiated indicates the merchant initiated the cancellation.
	CancelReasonMerchantInitiated CancelReason = 2

	// CancelReasonTransactionNotFound indicates the transaction was not found during a cancel attempt.
	CancelReasonTransactionNotFound CancelReason = 3

	// CancelReasonTimeout indicates the transaction was cancelled due to timeout (12 hours).
	CancelReasonTimeout CancelReason = 4
)

// ReceiptState represents the state of a receipt in the Subscribe API.
type ReceiptState int

const (
	// ReceiptStateCreated indicates a newly created receipt (not yet paid).
	ReceiptStateCreated ReceiptState = 0

	// ReceiptStateSent indicates the receipt has been sent via SMS.
	ReceiptStateSent ReceiptState = 1

	// ReceiptStatePending indicates the receipt is pending.
	ReceiptStatePending ReceiptState = 2

	// ReceiptStatePaid indicates the receipt has been paid.
	ReceiptStatePaid ReceiptState = 4

	// ReceiptStateCancelled indicates the receipt has been cancelled.
	ReceiptStateCancelled ReceiptState = 20

	// ReceiptStateCancelledWithError indicates the receipt was cancelled with an error.
	ReceiptStateCancelledWithError ReceiptState = 21
)
