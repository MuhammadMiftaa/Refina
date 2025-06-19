package dto

type ViewUserWallets struct {
	ID            string  `json:"id"`
	UserID        string  `json:"user_id"`
	WalletNumber  string  `json:"wallet_number"`
	WalletBalance float64 `json:"wallet_balance"`
	WalletName    string  `json:"wallet_name"`
	WalletTypeName string `json:"wallet_type_name"`
	WalletType    string  `json:"wallet_type"`
}