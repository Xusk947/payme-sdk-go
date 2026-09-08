package merchant

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/xusk947/payme-sdk-go/payme/rpc"
)

// contractResponse is a minimal JSON-RPC 2.0 response for contract assertions.
type contractResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *rpc.Error      `json:"error,omitempty"`
	ID      any             `json:"id"`
}

// parseContractResponse reads the recorder body and asserts HTTP 200 + JSON envelope.
func parseContractResponse(t *testing.T, w *httptest.ResponseRecorder) contractResponse {
	t.Helper()
	if w.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want 200 for all Merchant API responses", w.Code)
	}
	ct := w.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}
	var resp contractResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response body %q: %v", w.Body.String(), err)
	}
	if resp.JSONRPC != "2.0" {
		t.Errorf("jsonrpc = %q, want \"2.0\"", resp.JSONRPC)
	}
	// Exactly one of result or error must be present.
	hasResult := len(resp.Result) > 0 && string(resp.Result) != "null"
	hasError := resp.Error != nil
	if hasResult && hasError {
		t.Error("response has both result and error, expected exactly one")
	}
	if !hasResult && !hasError {
		t.Error("response has neither result nor error")
	}
	return resp
}

// requireRPCError asserts the response contains an error with the given code.
func requireRPCError(t *testing.T, resp contractResponse, wantCode int) {
	t.Helper()
	if resp.Error == nil {
		t.Fatalf("expected error code %d, got nil error", wantCode)
	}
	if resp.Error.Code != wantCode {
		t.Fatalf("error code = %d, want %d", resp.Error.Code, wantCode)
	}
}

// requireNullID asserts the response id is null (JSON null).
func requireNullID(t *testing.T, resp contractResponse) {
	t.Helper()
	if resp.ID != nil {
		t.Errorf("response id = %v, want null", resp.ID)
	}
}

// requireID asserts the response id matches the expected value.
func requireID(t *testing.T, resp contractResponse, want any) {
	t.Helper()
	if resp.ID != want {
		t.Errorf("response id = %v, want %v", resp.ID, want)
	}
}

// TestContract_AuthAllResponsesHTTP200 verifies that every Merchant API response
// is HTTP 200 with a JSON-RPC envelope, per the Payme protocol contract.
func TestContract_AuthAllResponsesHTTP200(t *testing.T) {
	h := NewHandler(&mockMerchantHandler{}, "login", "pass")

	tests := []struct {
		name    string
		method  string
		body    string
		authHdr string
	}{
		{
			name:    "missing auth header",
			method:  http.MethodPost,
			body:    `{"method":"CheckPerformTransaction","params":{},"id":1}`,
			authHdr: "",
		},
		{
			name:    "wrong password",
			method:  http.MethodPost,
			body:    `{"method":"CheckPerformTransaction","params":{},"id":1}`,
			authHdr: basicAuthHeader("login", "wrong"),
		},
		{
			name:    "wrong login",
			method:  http.MethodPost,
			body:    `{"method":"CheckPerformTransaction","params":{},"id":1}`,
			authHdr: basicAuthHeader("wrong", "pass"),
		},
		{
			name:    "bearer instead of basic",
			method:  http.MethodPost,
			body:    `{"method":"CheckPerformTransaction","params":{},"id":1}`,
			authHdr: "Bearer token123",
		},
		{
			name:    "damaged base64",
			method:  http.MethodPost,
			body:    `{"method":"CheckPerformTransaction","params":{},"id":1}`,
			authHdr: "Basic !!!invalid!!!",
		},
		{
			name:    "basic without colon",
			method:  http.MethodPost,
			body:    `{"method":"CheckPerformTransaction","params":{},"id":1}`,
			authHdr: "Basic " + base64.StdEncoding.EncodeToString([]byte("nocolon")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			if tt.authHdr != "" {
				req.Header.Set("Authorization", tt.authHdr)
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			resp := parseContractResponse(t, w)
			requireRPCError(t, resp, rpc.ErrCodeInsufficientPrivileges)
			requireNullID(t, resp)
		})
	}
}

// TestContract_AuthFailureDoesNotRevealWhichPartIsWrong verifies that all auth
// failure variants produce identical error codes and null id, so an attacker
// cannot distinguish missing, wrong-login, or wrong-password cases.
func TestContract_AuthFailureDoesNotRevealWhichPartIsWrong(t *testing.T) {
	h := NewHandler(&mockMerchantHandler{}, "login", "pass")

	variants := []string{
		"", // no header
		"Basic " + base64.StdEncoding.EncodeToString([]byte("login:wrong")),
		"Basic " + base64.StdEncoding.EncodeToString([]byte("wrong:pass")),
		"Basic !!!invalid!!!",
		"Bearer token123",
		"Basic " + base64.StdEncoding.EncodeToString([]byte("nocolon")),
	}

	for _, authHdr := range variants {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"method":"CheckPerformTransaction","params":{},"id":1}`))
		if authHdr != "" {
			req.Header.Set("Authorization", authHdr)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := parseContractResponse(t, w)
		requireRPCError(t, resp, rpc.ErrCodeInsufficientPrivileges)
		requireNullID(t, resp)
	}
}

// TestContract_NonPostReturnsInvalidHTTPMethod verifies that GET, PUT, DELETE,
// and PATCH all return -32300 with HTTP 200 and null id.
func TestContract_NonPostReturnsInvalidHTTPMethod(t *testing.T) {
	h := NewHandler(&mockMerchantHandler{}, "login", "pass")

	methods := []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodHead,
		http.MethodOptions,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/", nil)
			req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			resp := parseContractResponse(t, w)
			requireRPCError(t, resp, rpc.ErrCodeInvalidHTTPMethod)
			requireNullID(t, resp)
		})
	}
}

// TestContract_MalformedJSONReturnsParseError verifies that invalid JSON
// returns -32700 with HTTP 200.
func TestContract_MalformedJSONReturnsParseError(t *testing.T) {
	h := NewHandler(&mockMerchantHandler{}, "login", "pass")

	tests := []struct {
		name string
		body string
	}{
		{"empty body", ""},
		{"truncated json", `{"method":"CheckPerform`},
		{"invalid json", "not json at all"},
		{"just opening brace", "{"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			resp := parseContractResponse(t, w)
			requireRPCError(t, resp, rpc.ErrCodeParseError)
			requireNullID(t, resp)
		})
	}
}

// TestContract_UnknownMethodReturnsMethodNotFound verifies that an unknown
// JSON-RPC method returns -32601 with the request id preserved.
func TestContract_UnknownMethodReturnsMethodNotFound(t *testing.T) {
	h := NewHandler(&mockMerchantHandler{}, "login", "pass")

	body := `{"method":"UnknownMethod","params":{},"id":42}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp := parseContractResponse(t, w)
	requireRPCError(t, resp, rpc.ErrCodeMethodNotFound)
	// id should be preserved for method-not-found (parse succeeded).
	if resp.ID == nil {
		t.Error("expected non-null id for method-not-found, got null")
	}
}

// TestContract_ValidRequestPreservesID verifies that a successful request
// returns the same id that was sent.
func TestContract_ValidRequestPreservesID(t *testing.T) {
	mock := &mockMerchantHandler{
		checkPerformResult: &CheckPerformTransactionResponse{Allow: true},
	}
	h := NewHandler(mock, "login", "pass")

	body := `{"method":"CheckPerformTransaction","params":{"id":"5305e3bab097f420a62ced0b","time":1399114284039,"amount":500000,"account":{"phone":"903595731"}},"id":"req-abc"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp := parseContractResponse(t, w)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}
	requireID(t, resp, "req-abc")
}

// TestContract_NoInternalDetailsInErrorMessages verifies that error messages
// do not leak SQL, stack traces, or internal implementation details.
func TestContract_NoInternalDetailsInErrorMessages(t *testing.T) {
	mock := &mockMerchantHandler{
		checkPerformErr: &customError{msg: "sql: no rows in result set (SELECT * FROM users WHERE id=$1)"},
	}
	h := NewHandler(mock, "login", "pass")

	body := `{"method":"CheckPerformTransaction","params":{"id":"5305e3bab097f420a62ced0b","time":1399114284039,"amount":500000,"account":{"phone":"903595731"}},"id":1}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Authorization", basicAuthHeader("login", "pass"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp := parseContractResponse(t, w)
	requireRPCError(t, resp, rpc.ErrCodeInternal)
	// The internal error message should be the generic "Internal error", not the
	// raw error string from the handler.
	if resp.Error.Message["en"] != "Internal error" {
		t.Errorf("error message = %q, want \"Internal error\" (no internal details leaked)", resp.Error.Message["en"])
	}
}

// TestContract_CaseInsensitiveBasicPrefix verifies that "basic " prefix matching
// is case-insensitive per RFC 7617.
func TestContract_CaseInsensitiveBasicPrefix(t *testing.T) {
	mock := &mockMerchantHandler{
		checkPerformResult: &CheckPerformTransactionResponse{Allow: true},
	}
	h := NewHandler(mock, "login", "pass")

	creds := base64.StdEncoding.EncodeToString([]byte("login:pass"))

	tests := []struct {
		name  string
		value string
		want  bool // true = should succeed (result), false = should fail (auth error)
	}{
		{"uppercase Basic", "Basic " + creds, true},
		{"lowercase basic", "basic " + creds, true},
		{"mixed case BAsIc", "BAsIc " + creds, true},
		{"no space after basic", "Basic" + creds, false},
		{"wrong scheme", "Bearer " + creds, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"method":"CheckPerformTransaction","params":{},"id":1}`))
			req.Header.Set("Authorization", tt.value)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			resp := parseContractResponse(t, w)
			if tt.want {
				if resp.Error != nil {
					t.Errorf("expected success, got error: %v", resp.Error)
				}
			} else {
				requireRPCError(t, resp, rpc.ErrCodeInsufficientPrivileges)
			}
		})
	}
}
