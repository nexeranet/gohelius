package helius

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

const (
	BaseURL           = "https://api-mainnet.helius-rpc.com"
	RateLimitInterval = 200 * time.Millisecond // ~5 requests per second
)

type Client struct {
	client  *http.Client
	apiKey  string
	BaseURL string
	Limiter *rate.Limiter
}

func New(apiKey string) *Client {
	return &Client{
		client: &http.Client{
			Timeout: 3 * time.Minute,
		},
		apiKey:  apiKey,
		BaseURL: BaseURL,
		Limiter: rate.NewLimiter(rate.Every(RateLimitInterval), 1),
	}
}

func (c *Client) call(ctx context.Context, path string, query url.Values, target any) error {
	if err := c.Limiter.Wait(ctx); err != nil {
		return err
	}
	if target == nil {
		return errors.New("target is nil")
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, path)
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return err
	}
	if query == nil {
		query = url.Values{}
	}
	query.Set("api-key", c.apiKey)
	parsedURL.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", response.StatusCode, string(body))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

// GetTransactions returns parsed transaction history for a Solana wallet.
// Use beforeSignature for pagination (pass the last signature from previous page).
func (c *Client) GetTransactions(ctx context.Context, address string, limit int, beforeSignature string) ([]Transaction, error) {
	query := url.Values{}
	if limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", limit))
	}
	if beforeSignature != "" {
		query.Set("before", beforeSignature)
	}

	var result []Transaction
	path := fmt.Sprintf("/v0/addresses/%s/transactions", address)
	if err := c.call(ctx, path, query, &result); err != nil {
		return nil, err
	}
	return result, nil
}
