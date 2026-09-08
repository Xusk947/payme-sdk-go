package merchant

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/xusk947/payme-sdk-go/payme/rpc"
)

// Handler is an HTTP handler that processes Payme Business Merchant API
// JSON-RPC 2.0 requests. It validates authentication and dispatches
// requests to the provided MerchantHandler implementation.
type Handler struct {
	handler   MerchantHandler
	authLogin string
	authPass  string
}

// NewHandler creates a new Merchant API HTTP handler.
//
// The authLogin and authPass parameters are the credentials that Payme Business
// uses to authenticate with your server via HTTP Basic Auth. They should match
// the credentials configured in your Payme Business dashboard.
func NewHandler(h MerchantHandler, authLogin, authPass string) *Handler {
	return &Handler{
		handler:   h,
		authLogin: authLogin,
		authPass:  authPass,
	}
}

// rpcRequest is the JSON-RPC 2.0 request envelope.
type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      any             `json:"id"`
}

// rpcResponse is the JSON-RPC 2.0 response envelope.
type rpcResponse struct {
	JSONRPC string     `json:"jsonrpc"`
	Result  any        `json:"result,omitempty"`
	Error   *rpc.Error `json:"error,omitempty"`
	ID      any        `json:"id"`
}

// ServeHTTP processes incoming Payme Business Merchant API requests.
// It validates Basic Auth, parses the JSON-RPC request, dispatches to the
// appropriate MerchantHandler method, and returns the JSON-RPC response.
//
// All responses are returned with HTTP status 200, per the Payme Merchant API
// specification. Errors are represented as JSON-RPC error objects.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Validate HTTP Basic Auth
	if !h.validateAuth(r) {
		writeResponse(w, rpcResponse{
			JSONRPC: "2.0",
			Error:   rpc.ErrInsufficientPrivileges(),
			ID:      nil,
		})
		return
	}

	// Only POST is allowed; per Payme Merchant API spec, non-POST returns -32300.
	if r.Method != http.MethodPost {
		writeResponse(w, rpcResponse{
			JSONRPC: "2.0",
			Error:   rpc.ErrInvalidHTTPMethod(),
			ID:      nil,
		})
		return
	}

	// Read and parse the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, rpcResponse{
			JSONRPC: "2.0",
			Error:   rpc.ErrParseError(),
			ID:      nil,
		})
		return
	}
	defer func() { _ = r.Body.Close() }()

	var req rpcRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeResponse(w, rpcResponse{
			JSONRPC: "2.0",
			Error:   rpc.ErrParseError(),
			ID:      nil,
		})
		return
	}

	resp := h.dispatch(r.Context(), &req)
	resp.JSONRPC = "2.0"
	resp.ID = req.ID

	writeResponse(w, resp)
}

// dispatch routes the JSON-RPC request to the appropriate handler method.
func (h *Handler) dispatch(ctx context.Context, req *rpcRequest) rpcResponse {
	switch req.Method {
	case "CheckPerformTransaction":
		var params CheckPerformTransactionRequest
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return rpcResponse{Error: rpc.ErrParseError()}
		}
		result, err := h.handler.CheckPerformTransaction(ctx, &params)
		return buildResponse(result, err)

	case "CreateTransaction":
		var params CreateTransactionRequest
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return rpcResponse{Error: rpc.ErrParseError()}
		}
		result, err := h.handler.CreateTransaction(ctx, &params)
		return buildResponse(result, err)

	case "PerformTransaction":
		var params PerformTransactionRequest
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return rpcResponse{Error: rpc.ErrParseError()}
		}
		result, err := h.handler.PerformTransaction(ctx, &params)
		return buildResponse(result, err)

	case "CancelTransaction":
		var params CancelTransactionRequest
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return rpcResponse{Error: rpc.ErrParseError()}
		}
		result, err := h.handler.CancelTransaction(ctx, &params)
		return buildResponse(result, err)

	case "CheckTransaction":
		var params CheckTransactionRequest
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return rpcResponse{Error: rpc.ErrParseError()}
		}
		result, err := h.handler.CheckTransaction(ctx, &params)
		return buildResponse(result, err)

	case "GetStatement":
		var params GetStatementRequest
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return rpcResponse{Error: rpc.ErrParseError()}
		}
		result, err := h.handler.GetStatement(ctx, &params)
		return buildResponse(result, err)

	default:
		return rpcResponse{Error: rpc.ErrMethodNotFound(req.Method)}
	}
}

// buildResponse creates a JSON-RPC response from a result and error.
// Non-RPC errors are wrapped with a generic "Internal error" message to avoid
// leaking internal implementation details (SQL, stack traces, file paths) to
// the Payme caller. The original error is available to the handler for logging.
func buildResponse(result any, err error) rpcResponse {
	if err != nil {
		if rpcErr, ok := err.(*rpc.Error); ok {
			return rpcResponse{Error: rpcErr}
		}
		return rpcResponse{Error: rpc.ErrInternal()}
	}
	return rpcResponse{Result: result}
}

// validateAuth checks the HTTP Basic Auth header against the configured credentials
// using constant-time comparison to prevent timing attacks. Missing, malformed,
// and incorrect credentials all produce the same result and timing profile.
func (h *Handler) validateAuth(r *http.Request) bool {
	login, pass, ok := decodeBasicAuth(r.Header.Get("Authorization"))
	if !ok {
		return false
	}

	loginMatch := subtle.ConstantTimeCompare([]byte(login), []byte(h.authLogin))
	passMatch := subtle.ConstantTimeCompare([]byte(pass), []byte(h.authPass))
	return loginMatch == 1 && passMatch == 1
}

// decodeBasicAuth extracts the login and password from an HTTP Basic Auth header.
// Returns ok=false if the header is missing, not Basic, malformed base64, or
// has no colon separator. All failure cases are indistinguishable to the caller.
func decodeBasicAuth(authHeader string) (login, pass string, ok bool) {
	const prefix = "Basic "
	if len(authHeader) < len(prefix) || !strings.EqualFold(authHeader[:len(prefix)], prefix) {
		return "", "", false
	}
	encoded := strings.TrimSpace(authHeader[len(prefix):])
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", false
	}
	credentials := string(decoded)
	idx := strings.IndexByte(credentials, ':')
	if idx < 0 {
		return "", "", false
	}
	return credentials[:idx], credentials[idx+1:], true
}

// writeResponse writes the JSON-RPC response with HTTP 200 status.
func writeResponse(w http.ResponseWriter, resp rpcResponse) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		_, _ = fmt.Fprintf(w, `{"jsonrpc":"2.0","error":{"code":%d,"message":{"en":"Internal error"}},"id":null}`, rpc.ErrCodeInternal)
	}
}
