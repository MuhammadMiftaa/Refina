package view

type ViewUserInvestments struct {
	ID                 string  `json:"id"`
	UserID             string  `json:"user_id"`
	InvestmentType     string  `json:"investment_type"`
	InvestmentName     string  `json:"investment_name"`
	InvestmentAmount   float64 `json:"investment_amount"`
	InvestmentQuantity float64 `json:"investment_quantity"`
	InvestmentUnit     string  `json:"investment_unit"`
	InvestmentDate     string  `json:"investment_date"`
}
