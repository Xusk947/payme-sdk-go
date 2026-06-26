package models

// Amount represents a monetary amount in tiyin (1/100 of Uzbek Som).
// All amounts in the Payme API are expressed in tiyin as positive integers.
type Amount int64

// Timestamp represents a Unix timestamp in milliseconds.
type Timestamp int64

// ID represents a Payme Business identifier: a 24-character hex string.
type ID string
