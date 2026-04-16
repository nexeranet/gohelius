package helius

// WebhookType represents the delivery format of a webhook.
type WebhookType string

const (
	WebhookTypeEnhanced       WebhookType = "enhanced"
	WebhookTypeRaw            WebhookType = "raw"
	WebhookTypeDiscord        WebhookType = "discord"
	WebhookTypeEnhancedDevnet WebhookType = "enhancedDevnet"
	WebhookTypeRawDevnet      WebhookType = "rawDevnet"
	WebhookTypeDiscordDevnet  WebhookType = "discordDevnet"
)

// TransactionType represents the type of transaction that triggers a webhook.
type TransactionType string

const (
	TransactionTypeAny                TransactionType = "ANY"
	TransactionTypeSwap               TransactionType = "SWAP"
	TransactionTypeTransfer           TransactionType = "TRANSFER"
	TransactionTypeBurn               TransactionType = "BURN"
	TransactionTypeMintTo             TransactionType = "MINT_TO"
	TransactionTypeNFTSale            TransactionType = "NFT_SALE"
	TransactionTypeNFTMint            TransactionType = "NFT_MINT"
	TransactionTypeNFTListing         TransactionType = "NFT_LISTING"
	TransactionTypeNFTCancelListing   TransactionType = "NFT_CANCEL_LISTING"
	TransactionTypeNFTBid             TransactionType = "NFT_BID"
	TransactionTypeNFTCancelBid       TransactionType = "NFT_CANCEL_BID"
	TransactionTypeCompressedNFTMint  TransactionType = "COMPRESSED_NFT_MINT"
	TransactionTypeStakeSOL           TransactionType = "STAKE_SOL"
	TransactionTypeUnstakeSOL         TransactionType = "UNSTAKE_SOL"
	TransactionTypeTokenMint          TransactionType = "TOKEN_MINT"
)

// Webhook represents a Helius webhook configuration.
type Webhook struct {
	WebhookID        string            `json:"webhookID"`
	Wallet           string            `json:"wallet"`
	WebhookURL       string            `json:"webhookURL"`
	TransactionTypes []TransactionType `json:"transactionTypes"`
	AccountAddresses []string          `json:"accountAddresses"`
	WebhookType      WebhookType       `json:"webhookType"`
	AuthHeader       string            `json:"authHeader,omitempty"`
	Active           bool              `json:"active"`
}

// CreateWebhookRequest is the request body for creating a webhook.
type CreateWebhookRequest struct {
	WebhookURL       string            `json:"webhookURL"`
	TransactionTypes []TransactionType `json:"transactionTypes,omitempty"`
	AccountAddresses []string          `json:"accountAddresses,omitempty"`
	WebhookType      WebhookType       `json:"webhookType,omitempty"`
	AuthHeader       string            `json:"authHeader,omitempty"`
	Encoding         string            `json:"encoding,omitempty"`
	TxnStatus        string            `json:"txnStatus,omitempty"`
}

// UpdateWebhookRequest is the request body for updating a webhook.
type UpdateWebhookRequest struct {
	WebhookURL       string            `json:"webhookURL,omitempty"`
	TransactionTypes []TransactionType `json:"transactionTypes,omitempty"`
	AccountAddresses []string          `json:"accountAddresses,omitempty"`
	WebhookType      WebhookType       `json:"webhookType,omitempty"`
	AuthHeader       string            `json:"authHeader,omitempty"`
	Encoding         string            `json:"encoding,omitempty"`
	TxnStatus        string            `json:"txnStatus,omitempty"`
}

// ToggleWebhookRequest is the request body for toggling a webhook.
type ToggleWebhookRequest struct {
	Active bool `json:"active"`
}
