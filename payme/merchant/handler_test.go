package merchant

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/xusk947/payme-sdk-go/payme/models"
	"github.com/xusk947/payme-sdk-go/payme/rpc"
)

// mockMerchantHandler is a test implementation of MerchantHandler.
type mockMerchantHandler struct {
	checkPerformResult *CheckPerformTransactionResponse
	checkPerformErr    error

	createResult *CreateTransactionResponse
	createErr    error

	performResult *PerformTransactionResponse
	performErr    error

	cancelResult *CancelTransactionResponse
	cancelErr    error

	checkResult *CheckTransactionResponse
	checkErr    error

	statementResult *GetStatementResponse
	statementErr    error
}

func (m *mockMerchantHandler) CheckPerformTransaction(ctx context.Context, req *CheckPerformTransactionRequest) (*CheckPerformTransactionResponse, error) {
	return m.checkPerformResult, m.checkPerformErr
}

func (m *mockMerchantHandler) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error) {
	return m.createResult, m.createErr
}

func (m *mockMerchantHandler) PerformTransaction(ctx context.Context, req *PerformTransactionRequest) (*PerformTransactionResponse, error) {
	return m.performResult, m.performErr
}

func (m *mockMerchantHandler) CancelTransaction(ctx context.Context, req *CancelTransactionRequest) (*CancelTransactionResponse, error) {
	return m.cancelResult, m.cancelErr
}

func (m *mockMerchantHandler) CheckTransaction(ctx context.Context, req *CheckTransactionRequest) (*CheckTransactionResponse, error) {
	return m.checkResult, m.checkErr
}

func (m *mockMerchantHandler) GetStatement(ctx context.Context, req *GetStatementRequest) (*GetStatementResponse, error) {
	return m.statementResult, m.statementErr
}

func basicAuthHeader(login, pass string) string {
	creds := login + ":" + pass
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(creds))
}

func TestHandler_AuthFailure(t *testing.T) {
	h := NewHandler(&mockMerchantHandler{}, "login", "pass")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{}"))
	req.Header.Set("Authorization", "Basic wrongcreds")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp rpcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if resp.Error.Code != rpc.ErrCodeInsufficientPrivileges {
		t.Errorf("error code = %d, want %d", resp.Error.Code, rpc.ErrCodeInsufficientPrivileges)
	}
}

func TestHandler_NoAuthHeader(t *testing.T) {
	h := NewHandler(&mockMerchantHandler{}, "login", "pass")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp rpcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error == nil || resp.Error.Code != rpc.ErrCodeInsufficientPrivileges {
		t.Errorf("expected insufficient privileges error, got %v", resp.Error)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h := NewHandler(&mockMerchantHandler{}, "login", "pass")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	var resp rpcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error == nil || resp.Error.Code != rpc.ErrCodeInvalidRequest {
		t.Errorf("expected invalid request error, got %v", resp.Error)
	}
}

func TestHandler_InvalidJSON(t *testing.T) {
	h := NewHandler(&mockMerchantHandler{}, "login", "pass")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid json"))
	req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	var resp rpcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error == nil || resp.Error.Code != rpc.ErrCodeParseError {
		t.Errorf("expected parse error, got %v", resp.Error)
	}
}

func TestHandler_CheckPerformTransaction(t *testing.T) {
	mock := &mockMerchantHandler{
		checkPerformResult: &CheckPerformTransactionResponse{
			Allow: true,
		},
	}
	h := NewHandler(mock, "login", "pass")

	body := `{"method":"CheckPerformTransaction","params":{"id":"5305e3bab097f420a62ced0b","time":1399114284039,"amount":500000,"account":{"phone":"903595731"}},"id":1}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp rpcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}

	result, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatalf("expected map result, got %T", resp.Result)
	}

	if result["allow"] != true {
		t.Errorf("allow = %v, want true", result["allow"])
	}
}

func TestHandler_CreateTransaction(t *testing.T) {
	mock := &mockMerchantHandler{
		createResult: &CreateTransactionResponse{
			CreateTime: 1399114284039,
			Transaction: "5123",
			State:       models.StateCreated,
		},
	}
	h := NewHandler(mock, "login", "pass")

	body := `{"method":"CreateTransaction","params":{"id":"5305e3bab097f420a62ced0b","time":1399114284039,"amount":500000,"account":{"phone":"903595731"}},"id":2}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	var resp rpcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}

	result, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatalf("expected map result, got %T", resp.Result)
	}

	if result["transaction"] != "5123" {
		t.Errorf("transaction = %v, want 5123", result["transaction"])
	}
}

func TestHandler_PerformTransaction(t *testing.T) {
	mock := &mockMerchantHandler{
		performResult: &PerformTransactionResponse{
			Transaction: "5123",
			PerformTime: 1399114284039,
			State:       models.StateCompleted,
		},
	}
	h := NewHandler(mock, "login", "pass")

	body := `{"method":"PerformTransaction","params":{"id":"5305e3bab097f420a62ced0b"},"id":3}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	var resp rpcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}
}

func TestHandler_CancelTransaction(t *testing.T) {
	mock := &mockMerchantHandler{
		cancelResult: &CancelTransactionResponse{
			Transaction: "5123",
			CancelTime:  1399114284039,
			State:       models.StateCancelled,
		},
	}
	h := NewHandler(mock, "login", "pass")

	body := `{"method":"CancelTransaction","params":{"id":"5305e3bab097f420a62ced0b","reason":1},"id":4}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	var resp rpcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}
}

func TestHandler_CheckTransaction(t *testing.T) {
	mock := &mockMerchantHandler{
		checkResult: &CheckTransactionResponse{
			CreateTime:  1399114284039,
			PerformTime: 1399114285002,
			CancelTime:  0,
			Transaction: "5123",
			State:       models.StateCompleted,
			Reason:      nil,
		},
	}
	h := NewHandler(mock, "login", "pass")

	body := `{"method":"CheckTransaction","params":{"id":"5305e3bab097f420a62ced0b"},"id":5}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	var resp rpcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}
}

func TestHandler_GetStatement(t *testing.T) {
	mock := &mockMerchantHandler{
		statementResult: &GetStatementResponse{
			Transactions: []models.Transaction{},
		},
	}
	h := NewHandler(mock, "login", "pass")

	body := `{"method":"GetStatement","params":{"from":1399114284039,"to":1399120284000},"id":6}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	var resp rpcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}
}

func TestHandler_UnknownMethod(t *testing.T) {
	h := NewHandler(&mockMerchantHandler{}, "login", "pass")

	body := `{"method":"UnknownMethod","params":{},"id":7}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	var resp rpcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error == nil || resp.Error.Code != rpc.ErrCodeMethodNotFound {
		t.Errorf("expected method not found error, got %v", resp.Error)
	}
}

func TestHandler_HandlerReturnsRPCError(t *testing.T) {
	mock := &mockMerchantHandler{
		checkPerformErr: ErrAccountNotFound("phone"),
	}
	h := NewHandler(mock, "login", "pass")

	body := `{"method":"CheckPerformTransaction","params":{"id":"5305e3bab097f420a62ced0b","time":1399114284039,"amount":500000,"account":{"phone":"903595731"}},"id":8}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	var resp rpcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if resp.Error.Code != ErrCodeAccountNotFound {
		t.Errorf("error code = %d, want %d", resp.Error.Code, ErrCodeAccountNotFound)
	}
}

func TestHandler_HandlerReturnsGenericError(t *testing.T) {
	mock := &mockMerchantHandler{
		checkPerformErr: &customError{msg: "something went wrong"},
	}
	h := NewHandler(mock, "login", "pass")

	body := `{"method":"CheckPerformTransaction","params":{"id":"5305e3bab097f420a62ced0b","time":1399114284039,"amount":500000,"account":{"phone":"903595731"}},"id":9}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	var resp rpcResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if resp.Error.Code != rpc.ErrCodeInternal {
		t.Errorf("error code = %d, want %d", resp.Error.Code, rpc.ErrCodeInternal)
	}
}

type customError struct{ msg string }

func (e *customError) Error() string { return e.msg }

func TestValidateAuth(t *testing.T) {
	h := &Handler{authLogin: "login", authPass: "pass"}

	tests := []struct {
		name   string
		header string
		want   bool
	}{
		{"valid", basicAuthHeader("login", "pass"), true},
		{"wrong password", basicAuthHeader("login", "wrong"), false},
		{"wrong login", basicAuthHeader("wrong", "pass"), false},
		{"no auth header", "", false},
		{"non-basic auth", "Bearer token123", false},
		{"invalid base64", "Basic !!!invalid!!!", false},
		{"no colon", "Basic " + base64.StdEncoding.EncodeToString([]byte("nocolon")), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			got := h.validateAuth(req)
			if got != tt.want {
				t.Errorf("validateAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}
