package payme

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// PaymentOption configures a payment URL or form.
type PaymentOption func(*paymentConfig)

type paymentConfig struct {
	callback        string
	callbackTimeout int
	lang            string
	description     string
	detail          string
	testMode        bool
}

// WithCallback sets the callback URL for after payment or cancellation.
// The URL can contain parameters that Payme replaces:
//   - :transaction - transaction ID or "null" if transaction creation failed
//   - :account.{field} - account field values
func WithCallback(callbackURL string) PaymentOption {
	return func(c *paymentConfig) {
		c.callback = callbackURL
	}
}

// WithCallbackTimeout sets the timeout in milliseconds before redirecting
// to the callback URL after a successful payment. Default is 15ms.
func WithCallbackTimeout(ms int) PaymentOption {
	return func(c *paymentConfig) {
		c.callbackTimeout = ms
	}
}

// WithLang sets the language for the payment page. Valid values: "ru", "uz", "en".
// Default is "ru".
func WithLang(lang string) PaymentOption {
	return func(c *paymentConfig) {
		c.lang = lang
	}
}

// WithDescription sets the payment description.
func WithDescription(desc string) PaymentOption {
	return func(c *paymentConfig) {
		c.description = desc
	}
}

// WithDetail sets the payment detail object (items, shipping, discount).
// The detail is JSON-encoded and base64-encoded automatically.
func WithDetail(detail interface{}) PaymentOption {
	return func(c *paymentConfig) {
		jsonBytes, err := json.Marshal(detail)
		if err == nil {
			c.detail = base64.StdEncoding.EncodeToString(jsonBytes)
		}
	}
}

// WithPaymentTestMode sets the payment to use the test/sandbox checkout URL.
func WithPaymentTestMode() PaymentOption {
	return func(c *paymentConfig) {
		c.testMode = true
	}
}

// GeneratePaymentURL generates a GET URL for the Payme checkout page.
//
// The format is <checkout_url>/base64(params) where params are key=value
// pairs separated by semicolons, per the official Payme documentation.
//
// The merchantID is your Payme Business merchant ID.
// The amount is the payment amount in tiyin (1/100 of UZS).
// The account is a key-value map of account fields (e.g., {"order_id": "123"}).
//
// Optional payment options can be passed to customize the payment page.
//
// Example:
//
//	url := payme.GeneratePaymentURL("merchantID", 500000,
//	    map[string]string{"order_id": "123"},
//	    payme.WithLang("en"),
//	    payme.WithCallback("https://myshop.uz/payme/:transaction"),
//	)
func GeneratePaymentURL(merchantID string, amount int64, account map[string]string, opts ...PaymentOption) string {
	cfg := &paymentConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	baseURL := CheckoutURLProd
	if cfg.testMode {
		baseURL = CheckoutURLTest
	}

	var parts []string
	parts = append(parts, fmt.Sprintf("m=%s", merchantID))
	parts = append(parts, fmt.Sprintf("a=%d", amount))

	for key, value := range account {
		parts = append(parts, fmt.Sprintf("ac.%s=%s", key, value))
	}

	if cfg.lang != "" {
		parts = append(parts, fmt.Sprintf("l=%s", cfg.lang))
	}
	if cfg.callback != "" {
		parts = append(parts, fmt.Sprintf("cr=%s", cfg.callback))
	}
	if cfg.callbackTimeout > 0 {
		parts = append(parts, fmt.Sprintf("ct=%d", cfg.callbackTimeout))
	}
	if cfg.description != "" {
		parts = append(parts, fmt.Sprintf("ds=%s", cfg.description))
	}
	if cfg.detail != "" {
		parts = append(parts, fmt.Sprintf("dt=%s", cfg.detail))
	}

	joined := strings.Join(parts, ";")
	encoded := base64.StdEncoding.EncodeToString([]byte(joined))

	return fmt.Sprintf("%s/%s", baseURL, encoded)
}

// GeneratePaymentHTMLForm generates an HTML form for submitting a payment
// to the Payme checkout page via POST method.
//
// The merchantID is your Payme Business merchant ID.
// The amount is the payment amount in tiyin (1/100 of UZS).
// The account is a key-value map of account fields (e.g., {"order_id": "123"}).
//
// Optional payment options can be passed to customize the payment page.
//
// The generated form includes a submit button labeled "Pay with Payme".
func GeneratePaymentHTMLForm(merchantID string, amount int64, account map[string]string, opts ...PaymentOption) string {
	cfg := &paymentConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	actionURL := CheckoutURLProd
	if cfg.testMode {
		actionURL = CheckoutURLTest
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`<form method="POST" action="%s">`, actionURL))
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf(`  <input type="hidden" name="merchant" value="%s"/>`, merchantID))
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf(`  <input type="hidden" name="amount" value="%d"/>`, amount))
	sb.WriteString("\n")

	for key, value := range account {
		sb.WriteString(fmt.Sprintf(`  <input type="hidden" name="account[%s]" value="%s"/>`, key, value))
		sb.WriteString("\n")
	}

	if cfg.lang != "" {
		sb.WriteString(fmt.Sprintf(`  <input type="hidden" name="lang" value="%s"/>`, cfg.lang))
		sb.WriteString("\n")
	}
	if cfg.callback != "" {
		sb.WriteString(fmt.Sprintf(`  <input type="hidden" name="callback" value="%s"/>`, cfg.callback))
		sb.WriteString("\n")
	}
	if cfg.callbackTimeout > 0 {
		sb.WriteString(fmt.Sprintf(`  <input type="hidden" name="callback_timeout" value="%d"/>`, cfg.callbackTimeout))
		sb.WriteString("\n")
	}
	if cfg.description != "" {
		sb.WriteString(fmt.Sprintf(`  <input type="hidden" name="description" value="%s"/>`, cfg.description))
		sb.WriteString("\n")
	}
	if cfg.detail != "" {
		sb.WriteString(fmt.Sprintf(`  <input type="hidden" name="detail" value="%s"/>`, cfg.detail))
		sb.WriteString("\n")
	}

	sb.WriteString(`  <button type="submit">Pay with <b>Payme</b></button>`)
	sb.WriteString("\n")
	sb.WriteString("</form>")

	return sb.String()
}
