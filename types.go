package helius

type Transaction struct {
	Description      string            `json:"description"`
	Type             string            `json:"type"`
	Source           string            `json:"source"`
	Fee              int64             `json:"fee"`
	FeePayer         string            `json:"feePayer"`
	Signature        string            `json:"signature"`
	Slot             int64             `json:"slot"`
	Timestamp        int64             `json:"timestamp"`
	TokenTransfers   []TokenTransfer   `json:"tokenTransfers"`
	NativeTransfers  []NativeTransfer  `json:"nativeTransfers"`
	TransactionError *string           `json:"transactionError"`
	Events           map[string]any    `json:"events"`
}

type TokenTransfer struct {
	FromUserAccount string  `json:"fromUserAccount"`
	ToUserAccount   string  `json:"toUserAccount"`
	FromTokenAccount string `json:"fromTokenAccount"`
	ToTokenAccount   string `json:"toTokenAccount"`
	TokenAmount     float64 `json:"tokenAmount"`
	Mint            string  `json:"mint"`
	TokenStandard   string  `json:"tokenStandard"`
}

type NativeTransfer struct {
	FromUserAccount string `json:"fromUserAccount"`
	ToUserAccount   string `json:"toUserAccount"`
	Amount          int64  `json:"amount"`
}
