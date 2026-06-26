package models

// Account represents the account object used in Payme transactions.
// It is a flexible key-value map because the fields are determined by
// the merchant's business logic (e.g., phone, login, order_id).
type Account map[string]string
