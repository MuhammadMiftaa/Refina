package entity

import "github.com/google/uuid"

type Wallets struct {
	Base
	UserID       uuid.UUID `gorm:"type:uuid;not null"`
	WalletTypeID uuid.UUID `gorm:"type:uuid;not null"`
	Name         string    `gorm:"type:varchar(50);not null"`
	Number       string    `gorm:"type:varchar(50);not null"`
	Balance      float64   `gorm:"type:decimal(18,2);not null"`

	User       Users       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	WalletType WalletTypes `gorm:"foreignKey:WalletTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type WalletsResponse struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	WalletTypeID string  `json:"wallet_type_id"`
	Name         string  `json:"name"`
	Number       string  `json:"number"`
	Balance      float64 `json:"balance"`
}

type WalletsRequest struct {
	UserID       string  `json:"user_id"`
	WalletTypeID string  `json:"wallet_type_id"`
	Name         string  `json:"name"`
	Number       string  `json:"number"`
	Balance      float64 `json:"balance"`
}
