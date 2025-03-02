package entity

import "github.com/google/uuid"

type Attachments struct {
	Base
	TransactionID uuid.UUID `gorm:"type:uuid;not null"`
	Image         string    `gorm:"type:text"`

	Transaction Transactions `gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type AttachmentsResponse struct {
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	Image         string `json:"image"`
}

type AttachmentsRequest struct {
	TransactionID string `json:"transaction_id"`
	Image         string `json:"image"`
}
