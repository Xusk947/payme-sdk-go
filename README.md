# Payme Business SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/xusk947/payme-sdk-go.svg)](https://pkg.go.dev/github.com/xusk947/payme-sdk-go)
[![CI](https://github.com/xusk947/payme-sdk-go/actions/workflows/ci.yml/badge.svg)](https://github.com/xusk947/payme-sdk-go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/xusk947/payme-sdk-go)](https://goreportcard.com/report/github.com/xusk947/payme-sdk-go)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A Go SDK for the [Payme Business](https://developer.help.paycom.uz/) payment platform. Covers the Merchant API (server-side handler), Subscribe API (client-side caller), and payment initialization helpers.

## Features

- Merchant API: JSON-RPC 2.0 HTTP handler for receiving Payme Business callbacks
- Subscribe API: HTTP client for calling Payme Business endpoints (cards, receipts, invoices)
- Payment initialization: generate payment URLs and HTML forms for the Payme checkout page
- No external dependencies, pure Go standard library
- All methods accept `context.Context`
- 88%+ test coverage across all packages
- Tagged releases with GoReleaser

## Installation

```bash
go get github.com/xusk947/payme-sdk-go
```

## Quick Start

### Subscribe API (Client)

Use the Subscribe API to create receipts, pay with cards, send invoices, and manage card tokens.

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/xusk947/payme-sdk-go/payme"
    "github.com/xusk947/payme-sdk-go/payme/subscribe"
)

func main() {
    // Create a client (use WithTestMode() for sandbox)
    client := payme.NewClient("your_merchant_id", "your_key", payme.WithTestMode())

    ctx := context.Background()

    // Create a receipt
    detail := &subscribe.Detail{
        ReceiptType: 0,
        Items: []subscribe.Item{
            {Title: "Product A", Price: 500000, Count: 1, Code: "00702001001000001", VatPercent: 15},
        },
    }
    receipt, err := client.ReceiptsCreate(ctx, 500000, map[string]string{"order_id": "123"}, detail)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Receipt created: %s\n", receipt.ID)

    // Pay the receipt with a card token
    paid, err := client.ReceiptsPay(ctx, receipt.ID, "card_token_here", &subscribe.Payer{
        Phone: "998901234567",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Receipt paid, state: %d\n", paid.State)
}
```

### Merchant API (Server Handler)

Implement the `MerchantHandler` interface and register the handler with your HTTP server.

```go
package main

import (
    "context"
    "log"
    "net/http"

    "github.com/xusk947/payme-sdk-go/payme/merchant"
    "github.com/xusk947/payme-sdk-go/payme/models"
    "github.com/xusk947/payme-sdk-go/payme/rpc"
)

type myHandler struct{}

func (h *myHandler) CheckPerformTransaction(ctx context.Context, req *merchant.CheckPerformTransactionRequest) (*merchant.CheckPerformTransactionResponse, error) {
    // Validate account and amount
    if req.Amount <= 0 {
        return nil, merchant.ErrInvalidAmount("amount")
    }
    return &merchant.CheckPerformTransactionResponse{Allow: true}, nil
}

func (h *myHandler) CreateTransaction(ctx context.Context, req *merchant.CreateTransactionRequest) (*merchant.CreateTransactionResponse, error) {
    // Store transaction in your database
    return &merchant.CreateTransactionResponse{
        CreateTime:  1399114284039,
        Transaction: "your_internal_tx_id",
        State:       models.StateCreated,
    }, nil
}

func (h *myHandler) PerformTransaction(ctx context.Context, req *merchant.PerformTransactionRequest) (*merchant.PerformTransactionResponse, error) {
    // Mark order as paid
    return &merchant.PerformTransactionResponse{
        Transaction: "your_internal_tx_id",
        PerformTime: 1399114284039,
        State:       models.StateCompleted,
    }, nil
}

func (h *myHandler) CancelTransaction(ctx context.Context, req *merchant.CancelTransactionRequest) (*merchant.CancelTransactionResponse, error) {
    // Cancel/refund the transaction
    return &merchant.CancelTransactionResponse{
        Transaction: "your_internal_tx_id",
        CancelTime:  1399114284039,
        State:       models.StateCancelled,
    }, nil
}

func (h *myHandler) CheckTransaction(ctx context.Context, req *merchant.CheckTransactionRequest) (*merchant.CheckTransactionResponse, error) {
    // Return current transaction state
    return &merchant.CheckTransactionResponse{
        CreateTime:  1399114284039,
        PerformTime: 0,
        CancelTime:  0,
        Transaction: "your_internal_tx_id",
        State:       models.StateCreated,
        Reason:      nil,
    }, nil
}

func (h *myHandler) GetStatement(ctx context.Context, req *merchant.GetStatementRequest) (*merchant.GetStatementResponse, error) {
    // Return transactions for the period
    return &merchant.GetStatementResponse{
        Transactions: []models.Transaction{},
    }, nil
}

func main() {
    h := &myHandler{}
    handler := merchant.NewHandler(h, "your_login", "your_password")

    http.Handle("/payme", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Payment Initialization

Generate a payment URL or HTML form for the Payme checkout page.

```go
url := payme.GeneratePaymentURL("your_merchant_id", 500000,
    map[string]string{"order_id": "123"},
    payme.WithLang("en"),
    payme.WithCallback("https://myshop.uz/payme/:transaction"),
    payme.WithPaymentTestMode(),
)
// Redirect user to this URL

// Or generate an HTML form
form := payme.GeneratePaymentHTMLForm("your_merchant_id", 500000,
    map[string]string{"order_id": "123"},
    payme.WithLang("en"),
)
// Render this form in your HTML page
```

## API Reference

### Subscribe API Methods

| Method | Description | Auth |
|--------|-------------|------|
| `CardsCreate` | Create a card token | Partial (merchantID only) |
| `CardsGetVerifyCode` | Request SMS verification code | Partial |
| `CardsVerify` | Verify card with SMS code | Partial |
| `CardsCheck` | Check card token validity | Partial (merchantID only) |
| `CardsRemove` | Remove a card token | Full |
| `ReceiptsCreate` | Create a receipt/invoice | Full |
| `ReceiptsPay` | Pay a receipt with card token | Full |
| `ReceiptsSend` | Send invoice via SMS | Full |
| `ReceiptsCancel` | Cancel a paid receipt | Full |
| `ReceiptsCheck` | Check receipt status | Full |
| `ReceiptsGet` | Get full receipt info | Full |
| `ReceiptsGetAll` | Get all receipts for a period | Full |

### Merchant API Methods

| Method | Description |
|--------|-------------|
| `CheckPerformTransaction` | Check if transaction can be created |
| `CreateTransaction` | Create a financial transaction |
| `PerformTransaction` | Complete a transaction |
| `CancelTransaction` | Cancel a transaction |
| `CheckTransaction` | Check transaction state |
| `GetStatement` | Get transactions for a period |

## Error Handling

All errors from the Payme API are returned as `*rpc.Error` (aliased as `payme.RPCError`), which implements the `error` interface.

```go
receipt, err := client.ReceiptsGet(ctx, "receipt_id")
if err != nil {
    if rpcErr, ok := err.(*payme.RPCError); ok {
        fmt.Printf("Payme error [%d]: %s\n", rpcErr.Code, rpcErr.Message["en"])
    }
}
```

Error constructors with localized messages (ru, uz, en) are available in the `merchant` and `subscribe` packages.

## Configuration

| Option | Description |
|--------|-------------|
| `WithTestMode()` | Use sandbox endpoint |
| `WithHTTPClient(c)` | Custom HTTP client |
| `WithTimeout(d)` | HTTP client timeout |

## Testing

### Test Data

The SDK includes test data constants for the Payme Business sandbox environment:

```go
// Test cards (use ONLY in sandbox)
for _, card := range payme.TestCards {
    fmt.Printf("%s %s %s\n", card.Type, card.Number, card.Expire)
}
// Uzcard 8600495473316478 0399
// Uzcard 8600069195406311 0399
// Humo   9860010101010101 0399

// SMS code for all test cards is always 666666
fmt.Println(payme.TestSMSCode) // 666666

// Test merchant ID from documentation
fmt.Println(payme.TestMerchantID) // 5e730e8e0b852a417aa49ceb
```

### Running Tests

```bash
go test -race -cover ./...
```

## Contributing

PRs welcome. Run `go test -race -cover ./...` before submitting.

## License

[MIT](LICENSE)
