package merchant

import "github.com/xusk947/payme-sdk-go/payme/rpc"

// Merchant API error codes.
const (
	// ErrCodeAccountNotFound indicates the account was not found.
	ErrCodeAccountNotFound = -31050

	// ErrCodeInvalidAmount indicates the amount is invalid or does not match the billed amount.
	ErrCodeInvalidAmount = -31001

	// ErrCodeInvalidAccount indicates the account is invalid.
	ErrCodeInvalidAccount = -31002

	// ErrCodeTransactionNotFound indicates the transaction was not found.
	ErrCodeTransactionNotFound = -31003

	// ErrCodeTransactionAlreadyExists indicates a transaction with this ID already exists.
	ErrCodeTransactionAlreadyExists = -31004

	// ErrCodeTransactionAlreadyPerformed indicates the transaction was already performed.
	ErrCodeTransactionAlreadyPerformed = -31008

	// ErrCodeTransactionAlreadyCancelled indicates the transaction was already cancelled.
	ErrCodeTransactionAlreadyCancelled = -31009

	// ErrCodeCannotCancelTransaction indicates the transaction cannot be cancelled.
	ErrCodeCannotCancelTransaction = -31007

	// ErrCodeCannotPerformTransaction indicates the transaction cannot be performed.
	ErrCodeCannotPerformTransaction = -31006

	// ErrCodeInsufficientPrivileges indicates the merchant has insufficient privileges.
	ErrCodeInsufficientPrivileges = -32401
)

// ErrAccountNotFound returns an account-not-found error.
// The data parameter is typically the name of the account field that was not found.
func ErrAccountNotFound(data string) *rpc.Error {
	return rpc.NewError(ErrCodeAccountNotFound, map[string]string{
		"ru": "Аккаунт не найден",
		"uz": "Akkaunt topilmadi",
		"en": "Account not found",
	}, data)
}

// ErrInvalidAmount returns an invalid-amount error.
func ErrInvalidAmount(data string) *rpc.Error {
	return rpc.NewError(ErrCodeInvalidAmount, map[string]string{
		"ru": "Неверная сумма",
		"uz": "Noto'g'ri summa",
		"en": "Invalid amount",
	}, data)
}

// ErrInvalidAccount returns an invalid-account error.
func ErrInvalidAccount(data string) *rpc.Error {
	return rpc.NewError(ErrCodeInvalidAccount, map[string]string{
		"ru": "Неверный аккаунт",
		"uz": "Noto'g'ri akkaunt",
		"en": "Invalid account",
	}, data)
}

// ErrTransactionNotFound returns a transaction-not-found error.
func ErrTransactionNotFound(data string) *rpc.Error {
	return rpc.NewError(ErrCodeTransactionNotFound, map[string]string{
		"ru": "Транзакция не найдена",
		"uz": "Tranzaksiya topilmadi",
		"en": "Transaction not found",
	}, data)
}

// ErrTransactionAlreadyExists returns a transaction-already-exists error.
func ErrTransactionAlreadyExists(data string) *rpc.Error {
	return rpc.NewError(ErrCodeTransactionAlreadyExists, map[string]string{
		"ru": "Транзакция уже существует",
		"uz": "Tranzaksiya allaqachon mavjud",
		"en": "Transaction already exists",
	}, data)
}

// ErrTransactionAlreadyPerformed returns a transaction-already-performed error.
func ErrTransactionAlreadyPerformed(data string) *rpc.Error {
	return rpc.NewError(ErrCodeTransactionAlreadyPerformed, map[string]string{
		"ru": "Транзакция уже выполнена",
		"uz": "Tranzaksiya allaqachon bajarilgan",
		"en": "Transaction already performed",
	}, data)
}

// ErrTransactionAlreadyCancelled returns a transaction-already-cancelled error.
func ErrTransactionAlreadyCancelled(data string) *rpc.Error {
	return rpc.NewError(ErrCodeTransactionAlreadyCancelled, map[string]string{
		"ru": "Транзакция уже отменена",
		"uz": "Tranzaksiya allaqachon bekor qilingan",
		"en": "Transaction already cancelled",
	}, data)
}

// ErrCannotCancelTransaction returns a cannot-cancel-transaction error.
func ErrCannotCancelTransaction(data string) *rpc.Error {
	return rpc.NewError(ErrCodeCannotCancelTransaction, map[string]string{
		"ru": "Невозможно отменить транзакцию",
		"uz": "Tranzaksiyani bekor qilib bo'lmadi",
		"en": "Cannot cancel transaction",
	}, data)
}

// ErrCannotPerformTransaction returns a cannot-perform-transaction error.
func ErrCannotPerformTransaction(data string) *rpc.Error {
	return rpc.NewError(ErrCodeCannotPerformTransaction, map[string]string{
		"ru": "Невозможно выполнить транзакцию",
		"uz": "Tranzaksiyani bajarib bo'lmadi",
		"en": "Cannot perform transaction",
	}, data)
}
