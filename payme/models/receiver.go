package models

// Receiver represents a payment receiver in a chained payment transaction.
type Receiver struct {
	// ID is the 24-character Payme Business identifier of the receiver.
	ID string `json:"id"`

	// Amount is the amount to be transferred to this receiver, in tiyin.
	Amount int64 `json:"amount"`
}
