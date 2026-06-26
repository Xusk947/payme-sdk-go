package main

import (
	"context"
	"fmt"
	"log"

	"github.com/xusk947/payme-sdk-go/payme"
	"github.com/xusk947/payme-sdk-go/payme/subscribe"
)

func main() {
	// Create a client in test mode
	client := payme.NewClient("your_merchant_id", "your_key", payme.WithTestMode())

	ctx := context.Background()

	// Step 1: Create a card token (client-side method, partial auth)
	// Using a test card from payme.TestCards and payme.TestSMSCode
	testCard := payme.TestCards[0]
	token, err := client.CardsCreate(ctx, testCard.Number, testCard.Expire, true, nil)
	if err != nil {
		log.Fatalf("CardsCreate failed: %v", err)
	}
	fmt.Printf("Card token created: %s\n", token)

	// Step 2: Request verification code
	verifyResult, err := client.CardsGetVerifyCode(ctx, token, "998901234567")
	if err != nil {
		log.Fatalf("CardsGetVerifyCode failed: %v", err)
	}
	fmt.Printf("Verification code sent: %v, phone: %s\n", verifyResult.Sent, verifyResult.Phone)

	// Step 3: Verify the card with the SMS code
	// Test SMS code is always 666666 in sandbox
	verifiedToken, err := client.CardsVerify(ctx, token, payme.TestSMSCode)
	if err != nil {
		log.Fatalf("CardsVerify failed: %v", err)
	}
	fmt.Printf("Card verified, token: %s\n", verifiedToken)

	// Step 4: Create a receipt
	detail := &subscribe.Detail{
		ReceiptType: 0,
		Items: []subscribe.Item{
			{
				Title:      "Product A",
				Price:      500000,
				Count:      1,
				Code:       "00702001001000001",
				VatPercent: 15,
			},
		},
	}
	receipt, err := client.ReceiptsCreate(ctx, 500000, map[string]string{"order_id": "123"}, detail)
	if err != nil {
		log.Fatalf("ReceiptsCreate failed: %v", err)
	}
	fmt.Printf("Receipt created: %s\n", receipt.ID)

	// Step 5: Pay the receipt with the verified card token
	payer := &subscribe.Payer{
		Phone: "998901234567",
	}
	paid, err := client.ReceiptsPay(ctx, receipt.ID, verifiedToken, payer)
	if err != nil {
		log.Fatalf("ReceiptsPay failed: %v", err)
	}
	fmt.Printf("Receipt paid! State: %d\n", paid.State)

	// Step 6: Check receipt status
	state, err := client.ReceiptsCheck(ctx, receipt.ID)
	if err != nil {
		log.Fatalf("ReceiptsCheck failed: %v", err)
	}
	fmt.Printf("Receipt state: %d\n", state)
}
