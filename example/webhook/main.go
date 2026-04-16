package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	helius "github.com/nexeranet/gohelius"
)

func main() {
	apiKey := os.Getenv("HELIUS_API_KEY")
	if apiKey == "" {
		log.Fatal("HELIUS_API_KEY is required")
	}
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("WEBHOOK_URL is required (public URL that Helius can reach, e.g. ngrok)")
	}
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}
	accountAddress := os.Getenv("ACCOUNT_ADDRESS")
	if accountAddress == "" {
		log.Fatal("ACCOUNT_ADDRESS is required (Solana address to monitor)")
	}

	client := helius.New(apiKey, helius.APIBaseURLDevnet)
	client.BaseURL = helius.APIBaseURLDevnet
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	webhookID, err := ensureWebhook(ctx, client, webhookURL, accountAddress)
	if err != nil {
		log.Fatalf("Failed to ensure webhook: %v", err)
	}
	log.Printf("Webhook ready: %s", webhookID)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleWebhook)
	srv := &http.Server{Addr: listenAddr, Handler: mux}

	go func() {
		log.Printf("Listening on %s", listenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}

// ensureWebhook finds an existing webhook matching the URL or creates a new one.
func ensureWebhook(ctx context.Context, client *helius.Client, webhookURL, accountAddress string) (string, error) {
	webhooks, err := client.GetAllWebhooks(ctx)
	if err != nil {
		return "", fmt.Errorf("get all webhooks: %w", err)
	}
	for _, wh := range webhooks {
		if wh.WebhookURL == webhookURL {
			log.Printf("Found existing webhook %s for %s", wh.WebhookID, webhookURL)
			if !wh.Active {
				if _, err := client.ToggleWebhook(ctx, wh.WebhookID, true); err != nil {
					return "", fmt.Errorf("toggle webhook: %w", err)
				}
				log.Printf("Re-activated webhook %s", wh.WebhookID)
			}
			return wh.WebhookID, nil
		}
	}

	wh, err := client.CreateWebhook(ctx, helius.CreateWebhookRequest{
		WebhookURL:       webhookURL,
		TransactionTypes: []helius.TransactionType{helius.TransactionTypeAny},
		AccountAddresses: []string{accountAddress},
		WebhookType:      helius.WebhookTypeEnhancedDevnet,
	})
	if err != nil {
		return "", fmt.Errorf("create webhook: %w", err)
	}
	log.Printf("Created webhook %s -> %s", wh.WebhookID, webhookURL)
	return wh.WebhookID, nil
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read body: %v", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var events []json.RawMessage
	if err := json.Unmarshal(body, &events); err != nil {
		log.Printf("Failed to parse webhook payload: %v", err)
		log.Printf("Raw body: %s", string(body))
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Printf("Received %d event(s)", len(events))
	for i, raw := range events {
		var tx helius.Transaction
		if err := json.Unmarshal(raw, &tx); err != nil {
			log.Printf("  event[%d]: failed to parse: %v", i, err)
			log.Printf("  event[%d] raw: %s", i, string(raw))
			continue
		}
		log.Printf("  event[%d]: type=%s sig=%s fee=%d desc=%s",
			i, tx.Type, tx.Signature, tx.Fee, tx.Description)
		if len(tx.NativeTransfers) > 0 {
			for _, nt := range tx.NativeTransfers {
				log.Printf("    native: %s -> %s amount=%d", nt.FromUserAccount, nt.ToUserAccount, nt.Amount)
			}
		}
		if len(tx.TokenTransfers) > 0 {
			for _, tt := range tx.TokenTransfers {
				log.Printf("    token: %s -> %s mint=%s amount=%f", tt.FromUserAccount, tt.ToUserAccount, tt.Mint, tt.TokenAmount)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}
