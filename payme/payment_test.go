package payme

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestGeneratePaymentURL_Basic(t *testing.T) {
	u := GeneratePaymentURL("merchant123", 500000, map[string]string{"order_id": "123"})

	if !strings.HasPrefix(u, "https://paycom.uz/") {
		t.Fatalf("URL should start with https://paycom.uz/, got %q", u)
	}

	encoded := strings.TrimPrefix(u, "https://paycom.uz/")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("failed to decode base64: %v", err)
	}

	params := string(decoded)
	if !strings.Contains(params, "m=merchant123") {
		t.Errorf("params should contain m=merchant123, got %q", params)
	}
	if !strings.Contains(params, "a=500000") {
		t.Errorf("params should contain a=500000, got %q", params)
	}
	if !strings.Contains(params, "ac.order_id=123") {
		t.Errorf("params should contain ac.order_id=123, got %q", params)
	}
}

func TestGeneratePaymentURL_TestMode(t *testing.T) {
	u := GeneratePaymentURL("merchant", 100000, map[string]string{"order_id": "1"}, WithPaymentTestMode())

	if !strings.HasPrefix(u, "https://test.paycom.uz/") {
		t.Errorf("URL should start with https://test.paycom.uz/, got %q", u)
	}
}

func TestGeneratePaymentURL_WithOptions(t *testing.T) {
	u := GeneratePaymentURL("merchant", 500000, map[string]string{"order_id": "123"},
		WithLang("en"),
		WithCallback("https://myshop.uz/payme/:transaction"),
		WithCallbackTimeout(5000),
		WithDescription("Test payment"),
	)

	encoded := strings.TrimPrefix(u, "https://paycom.uz/")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("failed to decode base64: %v", err)
	}

	params := string(decoded)
	if !strings.Contains(params, "l=en") {
		t.Errorf("params should contain l=en, got %q", params)
	}
	if !strings.Contains(params, "cr=https://myshop.uz/payme/:transaction") {
		t.Errorf("params should contain callback, got %q", params)
	}
	if !strings.Contains(params, "ct=5000") {
		t.Errorf("params should contain ct=5000, got %q", params)
	}
	if !strings.Contains(params, "ds=Test payment") {
		t.Errorf("params should contain description, got %q", params)
	}
}

func TestGeneratePaymentURL_WithDetail(t *testing.T) {
	detail := map[string]any{
		"receipt_type": 0,
		"items": []map[string]any{
			{"title": "Item 1", "price": 500000, "count": 1},
		},
	}

	u := GeneratePaymentURL("merchant", 500000, map[string]string{"order_id": "1"}, WithDetail(detail))

	encoded := strings.TrimPrefix(u, "https://paycom.uz/")
	decodedParams, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("failed to decode base64 params: %v", err)
	}

	paramsStr := string(decodedParams)
	if !strings.Contains(paramsStr, "dt=") {
		t.Fatal("params should contain dt= (detail)")
	}

	idx := strings.Index(paramsStr, "dt=")
	encodedDetail := paramsStr[idx+3:]
	if idx2 := strings.Index(encodedDetail, ";"); idx2 >= 0 {
		encodedDetail = encodedDetail[:idx2]
	}
	if encodedDetail == "" {
		t.Fatal("detail is empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(encodedDetail)
	if err != nil {
		t.Fatalf("failed to decode detail: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(decoded, &result); err != nil {
		t.Fatalf("failed to unmarshal detail: %v", err)
	}

	if result["receipt_type"] != float64(0) {
		t.Errorf("receipt_type = %v, want 0", result["receipt_type"])
	}
}

func TestGeneratePaymentHTMLForm_Basic(t *testing.T) {
	form := GeneratePaymentHTMLForm("merchant123", 500000, map[string]string{"order_id": "123"})

	if !strings.Contains(form, `method="POST"`) {
		t.Error("form should have POST method")
	}
	if !strings.Contains(form, `action="https://paycom.uz"`) {
		t.Error("form should have production action URL")
	}
	if !strings.Contains(form, `name="merchant" value="merchant123"`) {
		t.Error("form should contain merchant field")
	}
	if !strings.Contains(form, `name="amount" value="500000"`) {
		t.Error("form should contain amount field")
	}
	if !strings.Contains(form, `name="account[order_id]" value="123"`) {
		t.Error("form should contain account field")
	}
	if !strings.Contains(form, "Pay with") {
		t.Error("form should contain submit button")
	}
}

func TestGeneratePaymentHTMLForm_TestMode(t *testing.T) {
	form := GeneratePaymentHTMLForm("merchant", 100000, map[string]string{"order_id": "1"}, WithPaymentTestMode())

	if !strings.Contains(form, `action="https://test.paycom.uz"`) {
		t.Error("form should have test action URL")
	}
}

func TestGeneratePaymentHTMLForm_WithOptions(t *testing.T) {
	form := GeneratePaymentHTMLForm("merchant", 500000, map[string]string{"order_id": "1"},
		WithLang("uz"),
		WithCallback("https://myshop.uz/callback"),
		WithDescription("Test"),
	)

	if !strings.Contains(form, `name="lang" value="uz"`) {
		t.Error("form should contain lang field")
	}
	if !strings.Contains(form, `name="callback" value="https://myshop.uz/callback"`) {
		t.Error("form should contain callback field")
	}
	if !strings.Contains(form, `name="description" value="Test"`) {
		t.Error("form should contain description field")
	}
}

func TestGeneratePaymentHTMLForm_WithDetail(t *testing.T) {
	detail := map[string]any{"receipt_type": 0}
	form := GeneratePaymentHTMLForm("merchant", 500000, map[string]string{"order_id": "1"}, WithDetail(detail))

	if !strings.Contains(form, `name="detail"`) {
		t.Error("form should contain detail field")
	}
}
