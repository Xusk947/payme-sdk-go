package payme

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xusk947/payme-sdk-go/payme/subscribe"
)

func TestReceiptsCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xAuth := r.Header.Get("X-Auth")
		if xAuth != "merchant:key" {
			t.Errorf("X-Auth = %q, want merchant:key", xAuth)
		}

		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req["method"] != "receipts.create" {
			t.Errorf("method = %v, want receipts.create", req["method"])
		}

		params, ok := req["params"].(map[string]any)
		if !ok {
			t.Fatalf("failed to assert params type")
		}
		if params["amount"] != float64(500000) {
			t.Errorf("amount = %v, want 500000", params["amount"])
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"result": map[string]any{
				"receipt": map[string]any{
					"_id":         "62da73b0803aced907a52b46",
					"create_time": 1658483632482,
					"state":       0,
					"amount":      500000,
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

	receipt, err := c.ReceiptsCreate(context.Background(), 500000, map[string]string{"order_id": "123"}, nil, "")
	if err != nil {
		t.Fatalf("ReceiptsCreate failed: %v", err)
	}
	if receipt.ID != "62da73b0803aced907a52b46" {
		t.Errorf("ID = %q, want 62da73b0803aced907a52b46", receipt.ID)
	}
	if receipt.Amount != 500000 {
		t.Errorf("Amount = %d, want 500000", receipt.Amount)
	}
}

func TestReceiptsCreate_WithDetail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"result": map[string]any{
				"receipt": map[string]any{
					"_id":   "test_id",
					"state": 0,
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

	detail := &subscribe.Detail{
		ReceiptType: 0,
		Items: []subscribe.Item{
			{Title: "Item 1", Price: 500000, Count: 1, Code: "00702001001000001", VatPercent: 15},
		},
	}

	receipt, err := c.ReceiptsCreate(context.Background(), 500000, map[string]string{"order_id": "1"}, detail, "")
	if err != nil {
		t.Fatalf("ReceiptsCreate failed: %v", err)
	}
	if receipt.ID != "test_id" {
		t.Errorf("ID = %q, want test_id", receipt.ID)
	}
}

func TestReceiptsPay(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req["method"] != "receipts.pay" {
			t.Errorf("method = %v, want receipts.pay", req["method"])
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"result": map[string]any{
				"receipt": map[string]any{
					"_id":      "receipt_123",
					"state":    4,
					"pay_time": 1481113810265,
					"card": map[string]any{
						"number": "860006******6311",
						"expire": "0399",
					},
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

	payer := &subscribe.Payer{Phone: "998901234567"}
	receipt, err := c.ReceiptsPay(context.Background(), "receipt_123", "card_token_456", payer)
	if err != nil {
		t.Fatalf("ReceiptsPay failed: %v", err)
	}
	if receipt.State != 4 {
		t.Errorf("State = %d, want 4", receipt.State)
	}
	if receipt.Card == nil || receipt.Card.Number != "860006******6311" {
		t.Errorf("Card number mismatch")
	}
}

func TestReceiptsSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req["method"] != "receipts.send" {
			t.Errorf("method = %v, want receipts.send", req["method"])
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"result": map[string]any{
				"success": true,
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	c := NewClient("merchant", "key")
	c.baseURL = server.URL

	success, err := c.ReceiptsSend(context.Background(), "receipt_123", "998901234567")
	if err != nil {
		t.Fatalf("ReceiptsSend failed: %v", err)
	}
	if !success {
		t.Error("success = false, want true")
	}
}

func TestReceiptsCancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req["method"] != "receipts.cancel" {
			t.Errorf("method = %v, want receipts.cancel", req["method"])
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"result": map[string]any{
				"receipt": map[string]any{
					"_id":         "receipt_123",
					"state":       20,
					"cancel_time": 1481113810265,
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

	receipt, err := c.ReceiptsCancel(context.Background(), "receipt_123")
	if err != nil {
		t.Fatalf("ReceiptsCancel failed: %v", err)
	}
	if receipt.State != 20 {
		t.Errorf("State = %d, want 20", receipt.State)
	}
}

func TestReceiptsCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req["method"] != "receipts.check" {
			t.Errorf("method = %v, want receipts.check", req["method"])
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"result": map[string]any{
				"state": 4,
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	c := NewClient("merchant", "key")
	c.baseURL = server.URL

	state, err := c.ReceiptsCheck(context.Background(), "receipt_123")
	if err != nil {
		t.Fatalf("ReceiptsCheck failed: %v", err)
	}
	if state != 4 {
		t.Errorf("state = %d, want 4", state)
	}
}

func TestReceiptsGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req["method"] != "receipts.get" {
			t.Errorf("method = %v, want receipts.get", req["method"])
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"result": map[string]any{
				"receipt": map[string]any{
					"_id":    "receipt_123",
					"state":  4,
					"amount": 500000,
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

	receipt, err := c.ReceiptsGet(context.Background(), "receipt_123")
	if err != nil {
		t.Fatalf("ReceiptsGet failed: %v", err)
	}
	if receipt.Amount != 500000 {
		t.Errorf("Amount = %d, want 500000", receipt.Amount)
	}
}

func TestReceiptsGetAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req["method"] != "receipts.get_all" {
			t.Errorf("method = %v, want receipts.get_all", req["method"])
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"result": []any{
				map[string]any{
					"_id":   "receipt_1",
					"state": 4,
				},
				map[string]any{
					"_id":   "receipt_2",
					"state": 0,
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

	receipts, err := c.ReceiptsGetAll(context.Background(), 1399114284039, 1399120284000, 50, 0)
	if err != nil {
		t.Fatalf("ReceiptsGetAll failed: %v", err)
	}
	if len(receipts) != 2 {
		t.Errorf("receipts count = %d, want 2", len(receipts))
	}
	if receipts[0].ID != "receipt_1" {
		t.Errorf("first receipt ID = %q, want receipt_1", receipts[0].ID)
	}
}

func TestReceiptsCreate_Error(t *testing.T) {
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

	_, err := c.ReceiptsCreate(context.Background(), 500000, map[string]string{"order_id": "1"}, nil, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
