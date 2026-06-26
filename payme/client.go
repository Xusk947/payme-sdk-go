package payme

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

// Client is the HTTP client for the Payme Business Subscribe API.
// It sends JSON-RPC 2.0 requests to the Payme Business endpoints.
//
// Use NewClient to create a new instance with the appropriate configuration.
type Client struct {
	baseURL    string
	merchantID string
	key        string
	httpClient *http.Client
	idCounter  uint64
}

// Option configures a Client.
type Option func(*Client)

// WithTestMode sets the client to use the test/sandbox endpoint.
func WithTestMode() Option {
	return func(c *Client) {
		c.baseURL = BaseURLTest
	}
}

// WithHTTPClient sets a custom *http.Client for the API client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

// NewClient creates a new Subscribe API client.
//
// The merchantID is your Payme Business merchant ID (cashbox ID).
// The key is your Payme Business cashbox key/password.
//
// By default, the client uses the production endpoint. Use WithTestMode()
// to switch to the sandbox endpoint.
func NewClient(merchantID, key string, opts ...Option) *Client {
	c := &Client{
		baseURL:    BaseURLProd,
		merchantID: merchantID,
		key:        key,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Call sends a JSON-RPC 2.0 request to the Payme Business Subscribe API.
//
// The method parameter is the JSON-RPC method name (e.g., "receipts.create").
// The params parameter is the request parameters, which will be JSON-encoded.
//
// The authType parameter controls the X-Auth header format:
//   - "full": uses "merchantID:key" (for server-side methods like receipts.*)
//   - "partial": uses "merchantID" only (for client-side methods like cards.create)
//
// Returns the raw JSON result that the caller can unmarshal into the appropriate type.
func (c *Client) Call(ctx context.Context, method string, params any, authType string) (json.RawMessage, error) {
	id := atomic.AddUint64(&c.idCounter, 1)

	reqBody := map[string]any{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      id,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	switch authType {
	case "full":
		req.Header.Set("X-Auth", fmt.Sprintf("%s:%s", c.merchantID, c.key))
	case "partial":
		req.Header.Set("X-Auth", c.merchantID)
	default:
		req.Header.Set("X-Auth", fmt.Sprintf("%s:%s", c.merchantID, c.key))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rpcResp struct {
		JSONRPC string          `json:"jsonrpc"`
		Result  json.RawMessage `json:"result"`
		Error   *RPCError       `json:"error"`
		ID      any             `json:"id"`
	}

	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, rpcResp.Error
	}

	return rpcResp.Result, nil
}

// callWithFullAuth is a convenience wrapper for server-side methods that require
// full authentication (merchantID:key).
func (c *Client) callWithFullAuth(ctx context.Context, method string, params any, result any) error {
	raw, err := c.Call(ctx, method, params, "full")
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, result); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}
	return nil
}

// callWithPartialAuth is a convenience wrapper for client-side methods that require
// partial authentication (merchantID only).
func (c *Client) callWithPartialAuth(ctx context.Context, method string, params any, result any) error {
	raw, err := c.Call(ctx, method, params, "partial")
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, result); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}
	return nil
}
