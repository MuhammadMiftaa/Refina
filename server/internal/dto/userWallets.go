package dto

type WalletResponse struct {
	ID         string  `json:"id"`
	Number     string  `json:"number"`
	Balance    float64 `json:"balance"`
	Name       string  `json:"name"`
	WalletType string  `json:"wallet_type"`
	Type       string  `json:"type"`
}

type UserWalletsResponse struct {
	UserID  string           `json:"user_id"`
	Name    string           `json:"name"`
	Email   string           `json:"email"`
	Wallets []WalletResponse `json:"wallets"`
}
