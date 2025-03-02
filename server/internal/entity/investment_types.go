package entity

type InvestmentType string

const (
	Gold                 InvestmentType = "gold"
	Stocks               InvestmentType = "stocks"
	MutualFunds          InvestmentType = "mutual funds"
	Bonds                InvestmentType = "bonds"
	GovernmentSecurities InvestmentType = "government securities"
	Deposits             InvestmentType = "deposits"
	OthersInvestment     InvestmentType = "others"
)

type InvestmentTypes struct {
	Base
	Name        string         `gorm:"type:varchar(50);not null"`
	Type        InvestmentType `gorm:"type:varchar(50);not null"`
	Description string         `gorm:"type:text"`
}

type InvestmentTypesResponse struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Type        InvestmentType `json:"type"`
	Description string         `json:"description"`
}

type InvestmentTypesRequest struct {
	Name        string         `json:"name"`
	Type        InvestmentType `json:"type"`
	Description string         `json:"description"`
}