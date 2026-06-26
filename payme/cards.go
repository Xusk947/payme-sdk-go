package payme

import (
	"context"
	"fmt"
)

// CardsCreate creates a card token from card details.
// This is a client-side method: it uses partial auth (merchantID only).
//
// The cardNumber is the 16-digit card number without spaces.
// The cardExpire is the card expiry date in MMYY format (e.g., "0399").
// The save indicates whether to save the card for recurrent payments (optional, pass false if not needed).
// The account is optional account data to associate with the card.
//
// Returns the created card token.
func (c *Client) CardsCreate(ctx context.Context, cardNumber, cardExpire string, save bool, account map[string]string) (string, error) {
	params := map[string]any{
		"card": map[string]any{
			"number": cardNumber,
			"expire": cardExpire,
		},
		"save": save,
	}
	if account != nil {
		params["account"] = account
	}

	var result struct {
		Card struct {
			Token string `json:"token"`
		} `json:"card"`
	}

	if err := c.callWithPartialAuth(ctx, "cards.create", params, &result); err != nil {
		return "", err
	}

	return result.Card.Token, nil
}

// CardsGetVerifyCode requests a verification code (SMS) for a card token.
// This is a client-side method: it uses partial auth (merchantID only).
//
// The token is the card token returned by CardsCreate.
// The phone is the phone number associated with the card (optional, pass empty string if not needed).
//
// Returns the verify code result with sent status, masked phone, and wait time in milliseconds.
func (c *Client) CardsGetVerifyCode(ctx context.Context, token, phone string) (*GetVerifyCodeResult, error) {
	params := map[string]any{
		"token": token,
	}
	if phone != "" {
		params["phone"] = phone
	}

	var result GetVerifyCodeResult
	if err := c.callWithPartialAuth(ctx, "cards.get_verify_code", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CardsVerify verifies a card token using the SMS code sent to the cardholder.
// This is a client-side method: it uses partial auth (merchantID only).
//
// The token is the card token returned by CardsCreate.
// The code is the SMS verification code sent to the cardholder's phone.
//
// Returns the verified card token.
func (c *Client) CardsVerify(ctx context.Context, token, code string) (string, error) {
	params := map[string]any{
		"token": token,
		"code":  code,
	}

	var result struct {
		Card struct {
			Token string `json:"token"`
		} `json:"card"`
	}

	if err := c.callWithPartialAuth(ctx, "cards.verify", params, &result); err != nil {
		return "", err
	}

	return result.Card.Token, nil
}

// CardsCheck checks the validity of a card token.
// This is a client-side method: it uses partial auth (merchantID only).
//
// The token is the card token to check.
//
// Returns card information including masked number, expiry, and verification status.
func (c *Client) CardsCheck(ctx context.Context, token string) (*CardCheckResult, error) {
	params := map[string]any{
		"token": token,
	}

	var result struct {
		Card CardCheckResult `json:"card"`
	}

	if err := c.callWithPartialAuth(ctx, "cards.check", params, &result); err != nil {
		return nil, err
	}

	return &result.Card, nil
}

// CardsRemove removes a card token.
// This is a server-side method: it uses full auth (merchantID:key).
//
// The token is the card token to remove.
//
// Returns true if the card was successfully removed.
func (c *Client) CardsRemove(ctx context.Context, token string) (bool, error) {
	params := map[string]any{
		"token": token,
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := c.callWithFullAuth(ctx, "cards.remove", params, &result); err != nil {
		return false, err
	}

	return result.Success, nil
}

// CardCheckResult contains the result of a cards.check call.
type CardCheckResult struct {
	Token     string `json:"token"`
	Number    string `json:"number"`
	Expire    string `json:"expire"`
	Recurrent bool   `json:"recurrent"`
	Verify    bool   `json:"verify"`
}

// GetVerifyCodeResult contains the result of a cards.get_verify_code call.
type GetVerifyCodeResult struct {
	Sent  bool   `json:"sent"`
	Phone string `json:"phone"`
	Wait  int64  `json:"wait"`
}

// String returns a string representation of the card check result.
func (c *CardCheckResult) String() string {
	return fmt.Sprintf("Card{Number: %s, Expire: %s, Verified: %v, Recurrent: %v}",
		c.Number, c.Expire, c.Verify, c.Recurrent)
}
