package payme

import "testing"

func TestTestCards(t *testing.T) {
	if len(TestCards) == 0 {
		t.Fatal("TestCards is empty")
	}

	for i, card := range TestCards {
		if card.Number == "" {
			t.Errorf("TestCards[%d]: number is empty", i)
		}
		if len(card.Number) != 16 {
			t.Errorf("TestCards[%d]: number length = %d, want 16", i, len(card.Number))
		}
		if card.Expire == "" {
			t.Errorf("TestCards[%d]: expire is empty", i)
		}
		if card.Type == "" {
			t.Errorf("TestCards[%d]: type is empty", i)
		}
	}
}

func TestTestSMSCode(t *testing.T) {
	if TestSMSCode != "666666" {
		t.Errorf("TestSMSCode = %q, want 666666", TestSMSCode)
	}
}

func TestTestMerchantID(t *testing.T) {
	if TestMerchantID == "" {
		t.Error("TestMerchantID is empty")
	}
}
