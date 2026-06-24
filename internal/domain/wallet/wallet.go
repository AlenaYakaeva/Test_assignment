package wallet

type Wallet struct {
	WalletID string
	Amount   float64
}

type WalletRequest struct {
	WalletID      string  `json:"walletID" validate:"required"`
	OperationType string  `json:"operationType"`
	Amount        float64 `json:"amount"`
}
