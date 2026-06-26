package payme

// Test data for the Payme Business sandbox environment.
//
// These constants are sourced from the official Payme Business documentation
// (https://developer.help.paycom.uz/) and can be used for testing in the
// sandbox environment only. Do NOT use these in production.
const (
	// TestMerchantID is a test merchant ID used in Payme Business documentation examples.
	TestMerchantID = "5e730e8e0b852a417aa49ceb"

	// TestSMSCode is the SMS verification code for all test cards in the sandbox.
	// The code is always 666666 for every test card.
	TestSMSCode = "666666"
)

// TestCard represents a test card for the Payme Business sandbox environment.
//
// Use these cards ONLY in the test/sandbox environment. The SMS verification
// code for all test cards is always TestSMSCode ("666666").
type TestCard struct {
	// Number is the 16-digit card number without spaces.
	Number string

	// Expire is the card expiry date in MMYY format (e.g., "0399").
	Expire string

	// Type is the card payment system ("Uzcard" or "Humo").
	Type string
}

// TestCards is the list of test cards provided by Payme Business for sandbox testing.
//
// Source: https://developer.help.paycom.uz/protokol-subscribe-api
var TestCards = []TestCard{
	{
		Number: "8600495473316478",
		Expire: "0399",
		Type:   "Uzcard",
	},
	{
		Number: "8600069195406311",
		Expire: "0399",
		Type:   "Uzcard",
	},
	{
		Number: "9860010101010101",
		Expire: "0399",
		Type:   "Humo",
	},
}
