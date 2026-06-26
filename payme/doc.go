// Package payme provides a Go SDK for the Payme Business payment platform.
//
// The SDK supports two protocols:
//
// # Merchant API
//
// The Merchant API is a server-side JSON-RPC 2.0 handler that receives
// requests from Payme Business. You implement the MerchantHandler interface
// and register the handler with your HTTP server.
//
// Example:
//
//	handler := payme.NewMerchantHandler(myHandler, "myLogin", "myPassword")
//	http.Handle("/payme", handler)
//	http.ListenAndServe(":8080", nil)
//
// # Subscribe API
//
// The Subscribe API is a client that sends JSON-RPC 2.0 requests to Payme
// Business endpoints. Use it to create receipts, pay with cards, send invoices,
// and manage card tokens.
//
// Example:
//
//	client := payme.NewClient("merchantID", "key", payme.WithTestMode())
//	receipt, err := client.ReceiptsCreate(ctx, 500000, map[string]string{"order_id": "123"}, nil)
//
// # Payment Initialization
//
// Helpers are provided to generate payment URLs and HTML forms for the
// Payme checkout page.
//
//	url := payme.GeneratePaymentURL("merchantID", 500000, map[string]string{"order_id": "123"})
package payme
