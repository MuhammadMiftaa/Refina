package entity

import "github.com/google/uuid"

type CategoryType string

const (
	Income  CategoryType = "income"
	Expense CategoryType = "expense"
)

type Categories struct {
	Base
	ParentID *uuid.UUID   `gorm:"type:uuid"`
	Name     string       `gorm:"type:varchar(50);not null"`
	Type     CategoryType `gorm:"type:varchar(50);not null"`

	Parent   *Categories  `gorm:"foreignKey:ParentID;references:ID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"` // Relasi ke parent (opsional)
	Children []Categories `gorm:"foreignKey:ParentID"`                                                             // Relasi ke children (opsional)
}

type CategoriesResponse struct {
	ID       string       `json:"id"`
	ParentID string       `json:"parent_id"`
	Name     string       `json:"name"`
	Type     CategoryType `json:"type"`
}

type CategoriesRequest struct {
	ParentID string       `json:"parent_id"`
	Name     string       `json:"name"`
	Type     CategoryType `json:"type"`
}
