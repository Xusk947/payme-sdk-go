package subscribe

import "github.com/xusk947/payme-sdk-go/payme/rpc"

// Subscribe API error codes.
const (
	// ErrCodeCardNotFound indicates the card token was not found.
	ErrCodeCardNotFound = -31900

	// ErrCodeCardAlreadyExists indicates a card with this number already exists.
	ErrCodeCardAlreadyExists = -31901

	// ErrCodeCardExpired indicates the card has expired.
	ErrCodeCardExpired = -31902

	// ErrCodeCardInvalid indicates the card number is invalid.
	ErrCodeCardInvalid = -31903

	// ErrCodeCardTokenInvalid indicates the card token is invalid.
	ErrCodeCardTokenInvalid = -31904

	// ErrCodeVerifyCodeInvalid indicates the verification code is invalid.
	ErrCodeVerifyCodeInvalid = -31905

	// ErrCodeVerifyCodeExpired indicates the verification code has expired.
	ErrCodeVerifyCodeExpired = -31906

	// ErrCodeVerifyCodeAlreadySent indicates a verification code was already sent.
	ErrCodeVerifyCodeAlreadySent = -31907

	// ErrCodeReceiptNotFound indicates the receipt was not found.
	ErrCodeReceiptNotFound = -32001

	// ErrCodeReceiptAlreadyPaid indicates the receipt is already paid.
	ErrCodeReceiptAlreadyPaid = -32002

	// ErrCodeReceiptAlreadyCancelled indicates the receipt is already cancelled.
	ErrCodeReceiptAlreadyCancelled = -32003

	// ErrCodeInsufficientPrivileges indicates insufficient privileges for the operation.
	ErrCodeInsufficientPrivileges = -32401
)

// ErrCardNotFound returns a card-not-found error.
func ErrCardNotFound() *rpc.Error {
	return rpc.NewErrorSimple(ErrCodeCardNotFound, "Card not found", nil)
}

// ErrCardExpired returns a card-expired error.
func ErrCardExpired() *rpc.Error {
	return rpc.NewErrorSimple(ErrCodeCardExpired, "Card has expired", nil)
}

// ErrCardInvalid returns an invalid-card error.
func ErrCardInvalid() *rpc.Error {
	return rpc.NewErrorSimple(ErrCodeCardInvalid, "Invalid card number", nil)
}

// ErrCardTokenInvalid returns an invalid-card-token error.
func ErrCardTokenInvalid() *rpc.Error {
	return rpc.NewErrorSimple(ErrCodeCardTokenInvalid, "Invalid card token", nil)
}

// ErrVerifyCodeInvalid returns an invalid-verify-code error.
func ErrVerifyCodeInvalid() *rpc.Error {
	return rpc.NewErrorSimple(ErrCodeVerifyCodeInvalid, "Invalid verification code", nil)
}

// ErrVerifyCodeExpired returns an expired-verify-code error.
func ErrVerifyCodeExpired() *rpc.Error {
	return rpc.NewErrorSimple(ErrCodeVerifyCodeExpired, "Verification code has expired", nil)
}

// ErrReceiptNotFound returns a receipt-not-found error.
func ErrReceiptNotFound() *rpc.Error {
	return rpc.NewErrorSimple(ErrCodeReceiptNotFound, "Receipt not found", nil)
}

// ErrReceiptAlreadyPaid returns a receipt-already-paid error.
func ErrReceiptAlreadyPaid() *rpc.Error {
	return rpc.NewErrorSimple(ErrCodeReceiptAlreadyPaid, "Receipt is already paid", nil)
}

// ErrReceiptAlreadyCancelled returns a receipt-already-cancelled error.
func ErrReceiptAlreadyCancelled() *rpc.Error {
	return rpc.NewErrorSimple(ErrCodeReceiptAlreadyCancelled, "Receipt is already cancelled", nil)
}
