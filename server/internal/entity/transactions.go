package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transactions struct {
	Base
	WalletID        uuid.UUID `gorm:"type:uuid;not null"`
	CategoryID      uuid.UUID `gorm:"type:uuid;not null"`
	Amount          float64   `gorm:"type:decimal(18,2);not null"`
	TransactionDate time.Time `gorm:"type:timestamp;not null"`
	Description     string    `gorm:"type:text"`

	Wallet   Wallets    `gorm:"foreignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Category Categories `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

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
