package entity

import (
	"time"

	"github.com/google/uuid"
)

type Investments struct {
	Base
	InvestmentTypeID uuid.UUID `gorm:"type:uuid;not null"`
	UserID           uuid.UUID `gorm:"type:uuid;not null"`
	Name             string    `gorm:"type:varchar(50);not null"`
	Amount           float64   `gorm:"type:decimal(18,2);not null"`
	Quantity         float64   `gorm:"type:decimal(18,2);not null"`
	InvestmentDate   time.Time `gorm:"type:timestamp;not null"`
	Description      string    `gorm:"type:text"`

	InvestmentTypes InvestmentTypes `gorm:"foreignKey:InvestmentTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User            Users           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}