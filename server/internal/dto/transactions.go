package dto

import "time"

type TransactionsResponse struct {
	ID              string    `json:"id"`
	WalletID        string    `json:"wallet_id"`
	CategoryID      string    `json:"category_id"`
	Amount          float64   `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
	Description     string    `json:"description"`
}

type TransactionsRequest struct {
	WalletID        string    `json:"wallet_id"`
	CategoryID      string    `json:"category_id"`
	Amount          float64   `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
	Description     string    `json:"description"`
}
