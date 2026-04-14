package helius

import (
	"context"
	"os"
	"testing"
)

func getTestClient(t *testing.T) *Client {
	if testing.Short() {
		t.Skip("skipping live API test in short mode")
	}
	apiKey := os.Getenv("HELIUS_API_KEY")
	if apiKey == "" {
		t.Skip("HELIUS_API_KEY not set")
	}
	return New(apiKey)
}

func TestGetTransactions(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	txs, err := client.GetTransactions(ctx, "F8fGH4NMgZhzXRPE1YuG519HY1zpdHoczRx98Maka24e", 5, "")
	if err != nil {
		t.Fatalf("GetTransactions error: %v", err)
	}
	t.Logf("Got %d transactions", len(txs))
	for _, tx := range txs {
		t.Logf("  [%s] %s — %s", tx.Type, tx.Signature[:16], tx.Description)
	}
}
