package merchant

import (
	"testing"

	"github.com/xusk947/payme-sdk-go/payme/rpc"
)

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name string
		err  *rpc.Error
		code int
	}{
		{"ErrAccountNotFound", ErrAccountNotFound("phone"), ErrCodeAccountNotFound},
		{"ErrInvalidAmount", ErrInvalidAmount("amount"), ErrCodeInvalidAmount},
		{"ErrInvalidAccount", ErrInvalidAccount("account"), ErrCodeInvalidAccount},
		{"ErrTransactionNotFound", ErrTransactionNotFound("id"), ErrCodeTransactionNotFound},
		{"ErrTransactionAlreadyExists", ErrTransactionAlreadyExists("id"), ErrCodeTransactionAlreadyExists},
		{"ErrTransactionAlreadyPerformed", ErrTransactionAlreadyPerformed("id"), ErrCodeTransactionAlreadyPerformed},
		{"ErrTransactionAlreadyCancelled", ErrTransactionAlreadyCancelled("id"), ErrCodeTransactionAlreadyCancelled},
		{"ErrCannotCancelTransaction", ErrCannotCancelTransaction("id"), ErrCodeCannotCancelTransaction},
		{"ErrCannotPerformTransaction", ErrCannotPerformTransaction("id"), ErrCodeCannotPerformTransaction},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Errorf("Code = %d, want %d", tt.err.Code, tt.code)
			}
			if tt.err.Data == nil {
				t.Error("Data should not be nil")
			}
			if tt.err.Message["en"] == "" {
				t.Error("English message is empty")
			}
			if tt.err.Message["ru"] == "" {
				t.Error("Russian message is empty")
			}
			if tt.err.Message["uz"] == "" {
				t.Error("Uzbek message is empty")
			}
			if tt.err.Error() == "" {
				t.Error("Error() returned empty string")
			}
		})
	}
}
