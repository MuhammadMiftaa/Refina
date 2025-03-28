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