package helius

import (
	"context"
	"fmt"
	"net/http"
)

// CreateWebhook creates a new webhook.
func (c *Client) CreateWebhook(ctx context.Context, req CreateWebhookRequest) (*Webhook, error) {
	var result Webhook
	if err := c.callWithBody(ctx, http.MethodPost, "/v0/webhooks", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetWebhook retrieves a webhook by its ID.
func (c *Client) GetWebhook(ctx context.Context, webhookID string) (*Webhook, error) {
	var result Webhook
	path := fmt.Sprintf("/v0/webhooks/%s", webhookID)
	if err := c.call(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateWebhook updates an existing webhook.
func (c *Client) UpdateWebhook(ctx context.Context, webhookID string, req UpdateWebhookRequest) (*Webhook, error) {
	var result Webhook
	path := fmt.Sprintf("/v0/webhooks/%s", webhookID)
	if err := c.callWithBody(ctx, http.MethodPut, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ToggleWebhook activates or deactivates a webhook without removing it.
func (c *Client) ToggleWebhook(ctx context.Context, webhookID string, active bool) (*Webhook, error) {
	var result Webhook
	path := fmt.Sprintf("/v0/webhooks/%s", webhookID)
	if err := c.callWithBody(ctx, http.MethodPatch, path, ToggleWebhookRequest{Active: active}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteWebhook removes a webhook by its ID.
func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) error {
	path := fmt.Sprintf("/v0/webhooks/%s", webhookID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil, nil)
}

// GetAllWebhooks returns all webhooks for the account.
func (c *Client) GetAllWebhooks(ctx context.Context) ([]Webhook, error) {
	var result []Webhook
	if err := c.call(ctx, "/v0/webhooks", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}
