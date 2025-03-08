package dto

type InvestmentResponse struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	Quantity   float64 `json:"quantity"`
	Unit       string  `json:"unit"`
	InvestDate string  `json:"invest_date"`
}

type UserInvestmentsResponse struct {
	UserID      string               `json:"user_id"`
	Name        string               `json:"name"`
	Email       string               `json:"email"`
	Investments []InvestmentResponse `json:"investments"`
}
