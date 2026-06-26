package payme

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient_DefaultConfig(t *testing.T) {
	c := NewClient("merchant123", "key456")
	if c.merchantID != "merchant123" {
		t.Errorf("merchantID = %q, want %q", c.merchantID, "merchant123")
	}
	if c.key != "key456" {
		t.Errorf("key = %q, want %q", c.key, "key456")
	}
	if c.baseURL != BaseURLProd {
		t.Errorf("baseURL = %q, want %q", c.baseURL, BaseURLProd)
	}
	if c.httpClient.Timeout != 30*time.Second {
		t.Errorf("timeout = %v, want %v", c.httpClient.Timeout, 30*time.Second)
	}
}

func TestNewClient_TestMode(t *testing.T) {
	c := NewClient("merchant", "key", WithTestMode())
	if c.baseURL != BaseURLTest {
		t.Errorf("baseURL = %q, want %q", c.baseURL, BaseURLTest)
	}
}

func TestNewClient_CustomHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 5 * time.Second}
	c := NewClient("merchant", "key", WithHTTPClient(custom))
	if c.httpClient != custom {
		t.Error("custom HTTP client not set")
	}
}

func TestNewClient_WithTimeout(t *testing.T) {
	c := NewClient("merchant", "key", WithTimeout(10*time.Second))
	if c.httpClient.Timeout != 10*time.Second {
		t.Errorf("timeout = %v, want %v", c.httpClient.Timeout, 10*time.Second)
	}
}

func TestClient_Call_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}

		xAuth := r.Header.Get("X-Auth")
		expected := "merchant123:key456"
		if xAuth != expected {
			t.Errorf("X-Auth = %q, want %q", xAuth, expected)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", r.Header.Get("Content-Type"))
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"result": map[string]any{
				"receipt": map[string]any{
					"_id":   "test_receipt_id",
					"state": 0,
				},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	c := NewClient("merchant123", "key456")
	c.baseURL = server.URL

	result, err := c.Call(context.Background(), "receipts.create", map[string]any{
		"amount":  500000,
		"account": map[string]string{"order_id": "123"},
	}, "full")
	if err != nil {
		t.Fatalf("Call failed: %v", err)
	}

	if result == nil {
		t.Fatal("result is nil")
	}
}

func TestClient_Call_PartialAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xAuth := r.Header.Get("X-Auth")
		expected := "merchant123"
		if xAuth != expected {
			t.Errorf("X-Auth = %q, want %q", xAuth, expected)
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"result": map[string]any{
				"card": map[string]any{
					"token": "test_token",
				},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	c := NewClient("merchant123", "key456")
	c.baseURL = server.URL

	result, err := c.Call(context.Background(), "cards.create", map[string]any{
		"card": map[string]any{
			"number": "8600061234567890",
			"expire": "0399",
		},
	}, "partial")
	if err != nil {
		t.Fatalf("Call failed: %v", err)
	}

	if result == nil {
		t.Fatal("result is nil")
	}
}

func TestClient_Call_RPCError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"error": map[string]any{
				"code": -32001,
				"message": map[string]string{
					"en": "Receipt not found",
				},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	c := NewClient("merchant", "key")
	c.baseURL = server.URL

	_, err := c.Call(context.Background(), "receipts.get", map[string]any{
		"id": "nonexistent",
	}, "full")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	rpcErr, ok := err.(*RPCError)
	if !ok {
		t.Fatalf("expected *RPCError, got %T", err)
	}
	if rpcErr.Code != -32001 {
		t.Errorf("error code = %d, want -32001", rpcErr.Code)
	}
}

func TestClient_Call_Non200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := NewClient("merchant", "key")
	c.baseURL = server.URL

	_, err := c.Call(context.Background(), "receipts.create", map[string]any{}, "full")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClient_Call_InvalidJSONResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, "invalid json")
	}))
	defer server.Close()

	c := NewClient("merchant", "key")
	c.baseURL = server.URL

	_, err := c.Call(context.Background(), "receipts.create", map[string]any{}, "full")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClient_Call_ServerUnavailable(t *testing.T) {
	c := NewClient("merchant", "key")
	c.baseURL = "http://localhost:99999"

	_, err := c.Call(context.Background(), "receipts.create", map[string]any{}, "full")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClient_Call_DefaultAuthType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xAuth := r.Header.Get("X-Auth")
		expected := "merchant:key"
		if xAuth != expected {
			t.Errorf("X-Auth = %q, want %q", xAuth, expected)
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"result":  map[string]any{},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	c := NewClient("merchant", "key")
	c.baseURL = server.URL

	_, err := c.Call(context.Background(), "test.method", map[string]any{}, "unknown")
	if err != nil {
		t.Fatalf("Call failed: %v", err)
	}
}
