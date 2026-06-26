package subscribe

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
		{"ErrCardNotFound", ErrCardNotFound(), ErrCodeCardNotFound},
		{"ErrCardExpired", ErrCardExpired(), ErrCodeCardExpired},
		{"ErrCardInvalid", ErrCardInvalid(), ErrCodeCardInvalid},
		{"ErrCardTokenInvalid", ErrCardTokenInvalid(), ErrCodeCardTokenInvalid},
		{"ErrVerifyCodeInvalid", ErrVerifyCodeInvalid(), ErrCodeVerifyCodeInvalid},
		{"ErrVerifyCodeExpired", ErrVerifyCodeExpired(), ErrCodeVerifyCodeExpired},
		{"ErrReceiptNotFound", ErrReceiptNotFound(), ErrCodeReceiptNotFound},
		{"ErrReceiptAlreadyPaid", ErrReceiptAlreadyPaid(), ErrCodeReceiptAlreadyPaid},
		{"ErrReceiptAlreadyCancelled", ErrReceiptAlreadyCancelled(), ErrCodeReceiptAlreadyCancelled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Errorf("Code = %d, want %d", tt.err.Code, tt.code)
			}
			if tt.err.Message["en"] == "" {
				t.Error("English message is empty")
			}
			if tt.err.Error() == "" {
				t.Error("Error() returned empty string")
			}
		})
	}
}
