package rpc

import "testing"

func TestRPCError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "english message",
			err:  NewErrorSimple(-31050, "Account not found", "phone"),
			want: "payme error [-31050]: Account not found",
		},
		{
			name: "russian fallback",
			err: &Error{
				Code: -31050,
				Message: map[string]string{
					"ru": "Аккаунт не найден",
				},
			},
			want: "payme error [-31050]: Аккаунт не найден",
		},
		{
			name: "no message",
			err: &Error{
				Code: -32603,
			},
			want: "payme error [-32603]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	msg := map[string]string{"en": "test error"}
	err := NewError(-31050, msg, "field")
	if err.Code != -31050 {
		t.Errorf("Code = %d, want -31050", err.Code)
	}
	if err.Message["en"] != "test error" {
		t.Errorf("Message[en] = %q, want %q", err.Message["en"], "test error")
	}
	if err.Data != "field" {
		t.Errorf("Data = %v, want %q", err.Data, "field")
	}
}

func TestNewErrorSimple(t *testing.T) {
	err := NewErrorSimple(-32700, "Parse error", nil)
	if err.Code != -32700 {
		t.Errorf("Code = %d, want -32700", err.Code)
	}
	if err.Message["en"] != "Parse error" {
		t.Errorf("Message[en] = %q, want %q", err.Message["en"], "Parse error")
	}
	if err.Data != nil {
		t.Errorf("Data = %v, want nil", err.Data)
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		code int
	}{
		{"ErrParseError", ErrParseError(), ErrCodeParseError},
		{"ErrInvalidRequest", ErrInvalidRequest(), ErrCodeInvalidRequest},
		{"ErrInternal", ErrInternal(), ErrCodeInternal},
		{"ErrHTTPError", ErrHTTPError(), ErrCodeHTTPError},
		{"ErrInsufficientPrivileges", ErrInsufficientPrivileges(), ErrCodeInsufficientPrivileges},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Errorf("%s Code = %d, want %d", tt.name, tt.err.Code, tt.code)
			}
		})
	}
}

func TestErrMethodNotFound(t *testing.T) {
	err := ErrMethodNotFound("UnknownMethod")
	if err.Code != ErrCodeMethodNotFound {
		t.Errorf("Code = %d, want %d", err.Code, ErrCodeMethodNotFound)
	}
	if err.Data != "UnknownMethod" {
		t.Errorf("Data = %v, want %q", err.Data, "UnknownMethod")
	}
}

func TestErrInvalidParams(t *testing.T) {
	err := ErrInvalidParams("amount")
	if err.Code != ErrCodeInvalidParams {
		t.Errorf("Code = %d, want %d", err.Code, ErrCodeInvalidParams)
	}
	if err.Data != "amount" {
		t.Errorf("Data = %v, want %q", err.Data, "amount")
	}
}
