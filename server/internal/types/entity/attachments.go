package entity

import "github.com/google/uuid"

type Attachments struct {
	Base
	TransactionID uuid.UUID `gorm:"type:uuid;not null"`
	Image         string    `gorm:"type:text"`

	Transaction Transactions `gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}