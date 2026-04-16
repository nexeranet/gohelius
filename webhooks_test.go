package helius

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func getDevnetTestClient(t *testing.T) *Client {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping live API test in short mode")
	}
	client := getTestClient(t)
	client.BaseURL = APIBaseURLDevnet
	return client
}

func TestWebhookCRUD(t *testing.T) {
	client := getDevnetTestClient(t)
	ctx := context.Background()

	// Create
	wh, err := client.CreateWebhook(ctx, CreateWebhookRequest{
		WebhookURL:       "https://example.com/test-hook",
		TransactionTypes: []TransactionType{TransactionTypeAny},
		AccountAddresses: []string{"F8fGH4NMgZhzXRPE1YuG519HY1zpdHoczRx98Maka24e"},
		WebhookType:      WebhookTypeEnhancedDevnet,
	})
	if err != nil {
		t.Fatalf("CreateWebhook: %v", err)
	}
	t.Logf("Created webhook %s", wh.WebhookID)
	if wh.WebhookID == "" {
		t.Fatal("WebhookID is empty")
	}
	if wh.WebhookURL != "https://example.com/test-hook" {
		t.Fatalf("WebhookURL = %q", wh.WebhookURL)
	}

	webhookID := wh.WebhookID

	// Get
	got, err := client.GetWebhook(ctx, webhookID)
	if err != nil {
		t.Fatalf("GetWebhook: %v", err)
	}
	if got.WebhookID != webhookID {
		t.Fatalf("GetWebhook ID = %q, want %q", got.WebhookID, webhookID)
	}
	t.Logf("Got webhook %s, active=%v", got.WebhookID, got.Active)

	// Update
	updated, err := client.UpdateWebhook(ctx, webhookID, UpdateWebhookRequest{
		WebhookURL:       "https://example.com/updated-hook",
		TransactionTypes: []TransactionType{TransactionTypeTransfer},
		AccountAddresses: []string{"F8fGH4NMgZhzXRPE1YuG519HY1zpdHoczRx98Maka24e"},
		WebhookType:      WebhookTypeEnhancedDevnet,
	})
	if err != nil {
		t.Fatalf("UpdateWebhook: %v", err)
	}
	if updated.WebhookURL != "https://example.com/updated-hook" {
		t.Fatalf("UpdateWebhook URL = %q", updated.WebhookURL)
	}
	t.Logf("Updated webhook %s", updated.WebhookID)

	// Toggle off
	toggled, err := client.ToggleWebhook(ctx, webhookID, false)
	if err != nil {
		t.Fatalf("ToggleWebhook(false): %v", err)
	}
	if toggled.Active {
		t.Fatal("expected Active=false after toggle off")
	}
	t.Logf("Toggled webhook %s active=%v", toggled.WebhookID, toggled.Active)

	// Toggle on
	toggled, err = client.ToggleWebhook(ctx, webhookID, true)
	if err != nil {
		t.Fatalf("ToggleWebhook(true): %v", err)
	}
	if !toggled.Active {
		t.Fatal("expected Active=true after toggle on")
	}

	// Delete
	err = client.DeleteWebhook(ctx, webhookID)
	if err != nil {
		t.Fatalf("DeleteWebhook: %v", err)
	}
	t.Logf("Deleted webhook %s", webhookID)
}

func TestDeleteWebhook(t *testing.T) {
	client := getDevnetTestClient(t)
	ctx := context.Background()

	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	wh, err := client.CreateWebhook(ctx, CreateWebhookRequest{
		WebhookURL:       fmt.Sprintf("https://example.com/test-hook/%s", timeStamp),
		TransactionTypes: []TransactionType{TransactionTypeAny},
		AccountAddresses: []string{"F8fGH4NMgZhzXRPE1YuG519HY1zpdHoczRx98Maka24e"},
		WebhookType:      WebhookTypeEnhancedDevnet,
	})
	if err != nil {
		t.Fatalf("CreateWebhook: %v", err)
	}

	err = client.DeleteWebhook(ctx, wh.WebhookID)
	if err != nil {
		t.Fatalf("DeleteWebhook: %v", err)
	}
	t.Logf("Deleted webhook %s", wh.WebhookID)

	_, err = client.GetWebhook(ctx, wh.WebhookID)
	if err == nil {
		t.Fatal("expected error when getting deleted webhook")
	}
}

func TestGetAllWebhooks(t *testing.T) {
	client := getDevnetTestClient(t)
	ctx := context.Background()

	webhooks, err := client.GetAllWebhooks(ctx)
	if err != nil {
		t.Fatalf("GetAllWebhooks: %v", err)
	}
	t.Logf("Got %d webhooks", len(webhooks))
	for _, wh := range webhooks {
		t.Logf("  [%s] %s type=%s active=%v", wh.WebhookID, wh.WebhookURL, wh.WebhookType, wh.Active)
	}
}
