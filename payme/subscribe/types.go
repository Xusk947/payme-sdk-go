// Package subscribe defines types for the Payme Business Subscribe API,
// including receipts, cards, payers, and merchant information.
package subscribe

// Receipt represents a payment receipt in the Payme Subscribe API.
type Receipt struct {
	// ID is the unique receipt identifier in Payme Business.
	ID string `json:"_id"`

	// CreateTime is the timestamp in milliseconds when the receipt was created.
	CreateTime int64 `json:"create_time"`

	// PayTime is the timestamp in milliseconds when the receipt was paid.
	PayTime int64 `json:"pay_time"`

	// CancelTime is the timestamp in milliseconds when the receipt was cancelled.
	CancelTime int64 `json:"cancel_time"`

	// State is the current state of the receipt.
	State int `json:"state"`

	// Type is the receipt type.
	Type int `json:"type"`

	// External indicates whether the receipt is external.
	External bool `json:"external"`

	// Operation is the operation code.
	Operation int `json:"operation"`

	// Category is the receipt category (may be null).
	Category any `json:"category"`

	// Error contains error information if any.
	Error any `json:"error"`

	// Description is the payment description.
	Description string `json:"description"`

	// Detail contains the receipt detail (items, shipping, discount).
	Detail *Detail `json:"detail"`

	// Amount is the receipt amount in tiyin.
	Amount int64 `json:"amount"`

	// Currency is the currency code (860 = UZS).
	Currency int `json:"currency"`

	// Commission is the commission amount in tiyin.
	Commission int64 `json:"commission"`

	// Account is the list of account fields associated with the receipt.
	Account []AccountField `json:"account"`

	// Card contains card information if the receipt was paid by card.
	Card *CardInfo `json:"card"`

	// Merchant contains merchant information.
	Merchant *Merchant `json:"merchant"`

	// Meta contains metadata.
	Meta any `json:"meta"`

	// ProcessingID is the processing identifier.
	ProcessingID any `json:"processing_id"`
}

// Detail represents the detailed breakdown of a receipt, including items,
// shipping, and discount information.
type Detail struct {
	// ReceiptType is the fiscal receipt type (0 or 1).
	ReceiptType int `json:"receipt_type,omitempty"`

	// Shipping contains shipping/delivery information.
	Shipping *Shipping `json:"shipping,omitempty"`

	// Items is the list of items in the receipt.
	Items []Item `json:"items,omitempty"`

	// Discount contains discount information.
	Discount *Discount `json:"discount,omitempty"`
}

// Shipping represents shipping/delivery information in a receipt detail.
type Shipping struct {
	// Title is the shipping description.
	Title string `json:"title"`

	// Price is the shipping price in tiyin.
	Price int64 `json:"price"`
}

// Item represents a single item in a receipt.
type Item struct {
	// Discount is the discount amount in tiyin.
	Discount int64 `json:"discount,omitempty"`

	// Title is the item name/title.
	Title string `json:"title"`

	// Price is the price per unit in tiyin.
	Price int64 `json:"price"`

	// Count is the quantity of items.
	Count int `json:"count"`

	// Code is the IKPU (ИКПУ) code: identification code of products and services.
	Code string `json:"code,omitempty"`

	// Units is the unit code.
	Units int `json:"units,omitempty"`

	// VatPercent is the VAT percentage for this item.
	VatPercent int `json:"vat_percent,omitempty"`

	// PackageCode is the package code for the item.
	PackageCode string `json:"package_code,omitempty"`
}

// Discount represents discount information in a receipt detail.
type Discount struct {
	// Title is the discount description.
	Title string `json:"title"`

	// Price is the discount amount in tiyin.
	Price int64 `json:"price"`
}

// AccountField represents an account field associated with a receipt.
type AccountField struct {
	// Name is the field name.
	Name string `json:"name"`

	// Title is the display title of the field. Can be a string or a localized object
	// (e.g., {"ru": "Логин", "uz": "Login", "en": "Login"}).
	Title any `json:"title"`

	// Value is the field value.
	Value string `json:"value"`

	// Main indicates whether this is the main account field.
	Main bool `json:"main"`
}

// CardInfo represents card information in a paid receipt.
type CardInfo struct {
	// Number is the masked card number.
	Number string `json:"number"`

	// Expire is the card expiry date in MMYY format.
	Expire string `json:"expire"`
}

// Merchant represents merchant information in a receipt.
type Merchant struct {
	// ID is the merchant identifier.
	ID string `json:"_id"`

	// Name is the merchant name.
	Name string `json:"name"`

	// Organization is the organization name.
	Organization string `json:"organization"`

	// Address is the merchant address.
	Address string `json:"address"`

	// BusinessID is the business identifier.
	BusinessID string `json:"business_id,omitempty"`

	// EPOS contains the electronic point of sale information.
	EPOS *EPOS `json:"epos,omitempty"`

	// Date is the merchant registration date.
	Date int64 `json:"date"`

	// Logo is the merchant logo URL (may be null).
	Logo any `json:"logo"`

	// Type is the merchant type. Can be a string (e.g., "Shop") or a localized object
	// (e.g., {"ru": "Internet", "uz": "Internet"}).
	Type any `json:"type"`

	// Terms is the merchant terms (may be null).
	Terms any `json:"terms"`

	// Payer contains payer information (in paid receipts).
	Payer *Payer `json:"payer,omitempty"`
}

// EPOS represents electronic point of sale information.
type EPOS struct {
	// MerchantID is the EPOS merchant ID.
	MerchantID string `json:"merchantId"`

	// TerminalID is the EPOS terminal ID.
	TerminalID string `json:"terminalId"`
}

// Payer represents payer information for a receipt payment.
type Payer struct {
	// ID is the payer identifier (optional).
	ID string `json:"id,omitempty"`

	// Phone is the payer's phone number.
	Phone string `json:"phone,omitempty"`

	// Email is the payer's email address.
	Email string `json:"email,omitempty"`

	// Name is the payer's name.
	Name string `json:"name,omitempty"`

	// IP is the payer's IP address.
	IP string `json:"ip,omitempty"`
}

// Card represents a card token created via the Subscribe API.
type Card struct {
	// Token is the card token string.
	Token string `json:"token"`

	// Number is the masked card number.
	Number string `json:"number"`

	// Expire is the card expiry date in MMYY format.
	Expire string `json:"expire"`

	// Recurrent indicates whether the card supports recurrent payments.
	Recurrent bool `json:"recurrent"`

	// Verify indicates whether the card is verified.
	Verify bool `json:"verify"`
}
