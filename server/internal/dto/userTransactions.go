package dto

type RawTransactionsResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
}

type WalletWithTransactionsResponse struct {
	ID           string                 `json:"id"`
	Number       string                 `json:"number"`
	Balance      float64                `json:"balance"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	Transactions []RawTransactionsResponse `json:"transactions"`
}

type UserTransactionsResponse struct {
	UserID  string                           `json:"user_id"`
	Name    string                           `json:"name"`
	Email   string                           `json:"email"`
	Wallets []WalletWithTransactionsResponse `json:"wallets"`
}
