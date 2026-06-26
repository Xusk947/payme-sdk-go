// Command merchant-server is an example Payme Business Merchant API server.
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/xusk947/payme-sdk-go/payme/merchant"
	"github.com/xusk947/payme-sdk-go/payme/models"
)

// myHandler implements merchant.MerchantHandler.
// This is a minimal example: in production, use a persistent storage.
type myHandler struct {
	store map[string]*models.Transaction
}

func (h *myHandler) CheckPerformTransaction(_ context.Context, req *merchant.CheckPerformTransactionRequest) (*merchant.CheckPerformTransactionResponse, error) {
	if req.Amount <= 0 {
		return nil, merchant.ErrInvalidAmount("amount")
	}
	account := req.Account["order_id"]
	if account == "" {
		return nil, merchant.ErrAccountNotFound("order_id")
	}
	return &merchant.CheckPerformTransactionResponse{Allow: true}, nil
}

func (h *myHandler) CreateTransaction(_ context.Context, req *merchant.CreateTransactionRequest) (*merchant.CreateTransactionResponse, error) {
	if existing, ok := h.store[req.ID]; ok {
		return &merchant.CreateTransactionResponse{
			CreateTime:  existing.CreateTime,
			Transaction: existing.Transaction,
			State:       existing.State,
		}, nil
	}

	txID := "tx_" + req.ID
	h.store[req.ID] = &models.Transaction{
		ID:          req.ID,
		Amount:      req.Amount,
		Account:     req.Account,
		CreateTime:  time.Now().UnixMilli(),
		Transaction: txID,
		State:       models.StateCreated,
	}

	return &merchant.CreateTransactionResponse{
		CreateTime:  time.Now().UnixMilli(),
		Transaction: txID,
		State:       models.StateCreated,
	}, nil
}

func (h *myHandler) PerformTransaction(_ context.Context, req *merchant.PerformTransactionRequest) (*merchant.PerformTransactionResponse, error) {
	tx, ok := h.store[req.ID]
	if !ok {
		return nil, merchant.ErrTransactionNotFound("id")
	}
	if tx.State == models.StateCompleted {
		return &merchant.PerformTransactionResponse{
			Transaction: tx.Transaction,
			PerformTime: tx.PerformTime,
			State:       tx.State,
		}, nil
	}

	tx.State = models.StateCompleted
	tx.PerformTime = time.Now().UnixMilli()

	return &merchant.PerformTransactionResponse{
		Transaction: tx.Transaction,
		PerformTime: tx.PerformTime,
		State:       tx.State,
	}, nil
}

func (h *myHandler) CancelTransaction(_ context.Context, req *merchant.CancelTransactionRequest) (*merchant.CancelTransactionResponse, error) {
	tx, ok := h.store[req.ID]
	if !ok {
		return nil, merchant.ErrTransactionNotFound("id")
	}

	if tx.State == models.StateCancelled || tx.State == models.StateCancelledAfterComplete {
		return &merchant.CancelTransactionResponse{
			Transaction: tx.Transaction,
			CancelTime:  tx.CancelTime,
			State:       tx.State,
		}, nil
	}

	if tx.State == models.StateCompleted {
		tx.State = models.StateCancelledAfterComplete
	} else {
		tx.State = models.StateCancelled
	}
	tx.CancelTime = time.Now().UnixMilli()
	tx.Reason = (*models.CancelReason)(&req.Reason)

	return &merchant.CancelTransactionResponse{
		Transaction: tx.Transaction,
		CancelTime:  tx.CancelTime,
		State:       tx.State,
	}, nil
}

func (h *myHandler) CheckTransaction(_ context.Context, req *merchant.CheckTransactionRequest) (*merchant.CheckTransactionResponse, error) {
	tx, ok := h.store[req.ID]
	if !ok {
		return nil, merchant.ErrTransactionNotFound("id")
	}

	return &merchant.CheckTransactionResponse{
		CreateTime:  tx.CreateTime,
		PerformTime: tx.PerformTime,
		CancelTime:  tx.CancelTime,
		Transaction: tx.Transaction,
		State:       tx.State,
		Reason:      tx.Reason,
	}, nil
}

func (h *myHandler) GetStatement(_ context.Context, req *merchant.GetStatementRequest) (*merchant.GetStatementResponse, error) {
	var transactions []models.Transaction
	for _, tx := range h.store {
		if tx.CreateTime >= req.From && tx.CreateTime <= req.To {
			transactions = append(transactions, *tx)
		}
	}
	return &merchant.GetStatementResponse{Transactions: transactions}, nil
}

func main() {
	h := &myHandler{store: make(map[string]*models.Transaction)}

	handler := merchant.NewHandler(h, "your_login", "your_password")

	http.Handle("/payme", handler)
	log.Println("Merchant API server listening on :8080")

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
