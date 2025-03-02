package entity

type InvestmentTypes struct {
	Base
	Name string `gorm:"type:varchar(50);not null"`
	Unit string `gorm:"type:text"`
}

type InvestmentTypesResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type InvestmentTypesRequest struct {
	Name string `json:"name"`
	Unit string `json:"unit"`
}
