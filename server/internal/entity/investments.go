package entity

import (
	"time"

	"github.com/google/uuid"
)

type Investments struct {
	Base
	InvestmentsTypeID uuid.UUID `gorm:"type:uuid;not null"`
	UserID            uuid.UUID `gorm:"type:uuid;not null"`
	Name              string    `gorm:"type:varchar(50);not null"`
	Amount            float64   `gorm:"type:decimal(18,2);not null"`
	Quantity          float64   `gorm:"type:decimal(18,2);not null"`
	InvestmentDate    time.Time `gorm:"type:timestamp;not null"`
	Description       string    `gorm:"type:text"`

	InvestmentTypes InvestmentTypes `gorm:"foreignKey:InvestmentsTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User            Users           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type InvestmentsResponse struct {
	ID                string    `json:"id"`
	InvestmentsTypeID string    `json:"investments_type_id"`
	UserID            string    `json:"user_id"`
	Name              string    `json:"name"`
	Amount            float64   `json:"amount"`
	Quantity          float64   `json:"quantity"`
	InvestmentDate    time.Time `json:"investment_date"`
	Description       string    `json:"description"`
}

type InvestmentsRequest struct {
	InvestmentsTypeID string    `json:"investments_type_id"`
	UserID            string    `json:"user_id"`
	Name              string    `json:"name"`
	Amount            float64   `json:"amount"`
	Quantity          float64   `json:"quantity"`
	InvestmentDate    time.Time `json:"investment_date"`
	Description       string    `json:"description"`
}
