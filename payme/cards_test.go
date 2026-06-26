package payme

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCardsCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xAuth := r.Header.Get("X-Auth")
		if xAuth != "merchant123" {
			t.Errorf("X-Auth = %q, want merchant123 (partial auth)", xAuth)
		}

		var req map[string]any
		json.NewDecoder(r.Body).Decode(&req)

		if req["method"] != "cards.create" {
			t.Errorf("method = %v, want cards.create", req["method"])
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"result": map[string]any{
				"card": map[string]any{
					"token": "test_card_token_123",
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient("merchant123", "key456")
	c.baseURL = server.URL

	token, err := c.CardsCreate(context.Background(), "8600061234567890", "0399", true, nil)
	if err != nil {
		t.Fatalf("CardsCreate failed: %v", err)
	}
	if token != "test_card_token_123" {
		t.Errorf("token = %q, want test_card_token_123", token)
	}
}

func TestCardsGetVerifyCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		json.NewDecoder(r.Body).Decode(&req)

		if req["method"] != "cards.get_verify_code" {
			t.Errorf("method = %v, want cards.get_verify_code", req["method"])
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"result": map[string]any{
				"sent":  true,
				"phone": "99890*****31",
				"wait":  60000,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient("merchant", "key")
	c.baseURL = server.URL

	result, err := c.CardsGetVerifyCode(context.Background(), "token123", "998901234567")
	if err != nil {
		t.Fatalf("CardsGetVerifyCode failed: %v", err)
	}
	if !result.Sent {
		t.Error("sent = false, want true")
	}
	if result.Phone != "99890*****31" {
		t.Errorf("phone = %q, want 99890*****31", result.Phone)
	}
	if result.Wait != 60000 {
		t.Errorf("wait = %d, want 60000", result.Wait)
	}
}

func TestCardsVerify(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		json.NewDecoder(r.Body).Decode(&req)

		if req["method"] != "cards.verify" {
			t.Errorf("method = %v, want cards.verify", req["method"])
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"result": map[string]any{
				"card": map[string]any{
					"token": "verified_token_456",
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient("merchant", "key")
	c.baseURL = server.URL

	token, err := c.CardsVerify(context.Background(), "token123", "666666")
	if err != nil {
		t.Fatalf("CardsVerify failed: %v", err)
	}
	if token != "verified_token_456" {
		t.Errorf("token = %q, want verified_token_456", token)
	}
}

func TestCardsCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xAuth := r.Header.Get("X-Auth")
		expected := "merchant"
		if xAuth != expected {
			t.Errorf("X-Auth = %q, want %q (partial auth)", xAuth, expected)
		}

		var req map[string]any
		json.NewDecoder(r.Body).Decode(&req)

		if req["method"] != "cards.check" {
			t.Errorf("method = %v, want cards.check", req["method"])
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"result": map[string]any{
				"card": map[string]any{
					"token":     "token123",
					"number":    "860006******7890",
					"expire":    "0399",
					"recurrent": true,
					"verify":    true,
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient("merchant", "key")
	c.baseURL = server.URL

	card, err := c.CardsCheck(context.Background(), "token123")
	if err != nil {
		t.Fatalf("CardsCheck failed: %v", err)
	}
	if card.Number != "860006******7890" {
		t.Errorf("number = %q, want 860006******7890", card.Number)
	}
	if !card.Verify {
		t.Error("verify = false, want true")
	}
}

func TestCardsRemove(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xAuth := r.Header.Get("X-Auth")
		expected := "merchant:key"
		if xAuth != expected {
			t.Errorf("X-Auth = %q, want %q (full auth)", xAuth, expected)
		}

		var req map[string]any
		json.NewDecoder(r.Body).Decode(&req)

		if req["method"] != "cards.remove" {
			t.Errorf("method = %v, want cards.remove", req["method"])
		}

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"result": map[string]any{
				"success": true,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := NewClient("merchant", "key")
	c.baseURL = server.URL

	success, err := c.CardsRemove(context.Background(), "token123")
	if err != nil {
		t.Fatalf("CardsRemove failed: %v", err)
	}
	if !success {
		t.Error("success = false, want true")
	}
}
