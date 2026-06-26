package payme

import (
	"context"

	"github.com/xusk947/payme-sdk-go/payme/subscribe"
)

// ReceiptsCreate creates a new receipt (invoice) for payment.
// This is a server-side method: it uses full auth (merchantID:key).
//
// The amount is the payment amount in tiyin (1/100 of UZS).
// The account is a key-value map of account fields (e.g., {"order_id": "123"}).
// The detail is optional receipt detail (items, shipping, discount). Pass nil if not needed.
//
// Returns the created receipt.
func (c *Client) ReceiptsCreate(ctx context.Context, amount int64, account map[string]string, detail *subscribe.Detail) (*subscribe.Receipt, error) {
	params := map[string]any{
		"amount":  amount,
		"account": account,
	}
	if detail != nil {
		params["detail"] = detail
	}

	var result struct {
		Receipt subscribe.Receipt `json:"receipt"`
	}

	if err := c.callWithFullAuth(ctx, "receipts.create", params, &result); err != nil {
		return nil, err
	}

	return &result.Receipt, nil
}

// ReceiptsPay pays a receipt using a card token.
// This is a server-side method: it uses full auth (merchantID:key).
//
// The id is the receipt ID returned by ReceiptsCreate.
// The token is the verified card token.
// The payer contains payer information (phone, email, etc.).
//
// Returns the paid receipt.
func (c *Client) ReceiptsPay(ctx context.Context, id, token string, payer *subscribe.Payer) (*subscribe.Receipt, error) {
	params := map[string]any{
		"id":    id,
		"token": token,
	}
	if payer != nil {
		params["payer"] = payer
	}

	var result struct {
		Receipt subscribe.Receipt `json:"receipt"`
	}

	if err := c.callWithFullAuth(ctx, "receipts.pay", params, &result); err != nil {
		return nil, err
	}

	return &result.Receipt, nil
}

// ReceiptsSend sends an invoice (receipt) to a phone number via SMS.
// This is a server-side method: it uses full auth (merchantID:key).
//
// The id is the receipt ID.
// The phone is the recipient's phone number.
//
// Returns true if the SMS was sent successfully.
func (c *Client) ReceiptsSend(ctx context.Context, id, phone string) (bool, error) {
	params := map[string]any{
		"id":    id,
		"phone": phone,
	}

	var result struct {
		Success bool `json:"success"`
	}

	if err := c.callWithFullAuth(ctx, "receipts.send", params, &result); err != nil {
		return false, err
	}

	return result.Success, nil
}

// ReceiptsCancel cancels a paid receipt (puts it in the cancellation queue).
// This is a server-side method: it uses full auth (merchantID:key).
//
// The id is the receipt ID to cancel.
//
// Returns the cancelled receipt.
func (c *Client) ReceiptsCancel(ctx context.Context, id string) (*subscribe.Receipt, error) {
	params := map[string]any{
		"id": id,
	}

	var result struct {
		Receipt subscribe.Receipt `json:"receipt"`
	}

	if err := c.callWithFullAuth(ctx, "receipts.cancel", params, &result); err != nil {
		return nil, err
	}

	return &result.Receipt, nil
}

// ReceiptsCheck checks the status of a receipt.
// This is a server-side method: it uses full auth (merchantID:key).
//
// The id is the receipt ID to check.
//
// Returns the current state of the receipt.
// Possible state values: 0 (Created), 1 (Sent), 2 (Pending), 4 (Paid),
// 20 (Cancelled), 21 (Cancelled after paid).
func (c *Client) ReceiptsCheck(ctx context.Context, id string) (int, error) {
	params := map[string]any{
		"id": id,
	}

	var result struct {
		State int `json:"state"`
	}

	if err := c.callWithFullAuth(ctx, "receipts.check", params, &result); err != nil {
		return 0, err
	}

	return result.State, nil
}

// ReceiptsGet retrieves full information about a receipt.
// This is a server-side method: it uses full auth (merchantID:key).
//
// The id is the receipt ID.
//
// Returns the full receipt information.
func (c *Client) ReceiptsGet(ctx context.Context, id string) (*subscribe.Receipt, error) {
	params := map[string]any{
		"id": id,
	}

	var result struct {
		Receipt subscribe.Receipt `json:"receipt"`
	}

	if err := c.callWithFullAuth(ctx, "receipts.get", params, &result); err != nil {
		return nil, err
	}

	return &result.Receipt, nil
}

// ReceiptsGetAll retrieves all receipts for a specified time period.
// This is a server-side method: it uses full auth (merchantID:key).
//
// The from and to are timestamps in milliseconds defining the period.
// The count is the maximum number of receipts to return.
// The offset is the number of receipts to skip (for pagination).
//
// Returns the list of receipts.
func (c *Client) ReceiptsGetAll(ctx context.Context, from, to int64, count, offset int) ([]subscribe.Receipt, error) {
	params := map[string]any{
		"from":   from,
		"to":     to,
		"count":  count,
		"offset": offset,
	}

	var result []subscribe.Receipt

	if err := c.callWithFullAuth(ctx, "receipts.get_all", params, &result); err != nil {
		return nil, err
	}

	return result, nil
}
