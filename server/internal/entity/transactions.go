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

type TransactionsList struct {
	WalletID      string
	WalletNumber  string
	WalletBalance float64
	WalletName    string
	WalletType    string

	TransactionID   string
	CategoryName    string
	CategoryType    string
	Amount          float64
	TransactionDate string
	Description     string

	Image string
}
