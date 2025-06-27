package entity

type InvestmentTypes struct {
	Base
	Name string `gorm:"type:varchar(50);not null"`
	Unit string `gorm:"type:text"`
}