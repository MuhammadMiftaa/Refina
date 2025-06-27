package dto

import "time"

type InvestmentsResponse struct {
	ID               string    `json:"id"`
	InvestmentTypeID string    `json:"investment_type_id"`
	UserID           string    `json:"user_id"`
	Name             string    `json:"name"`
	Amount           float64   `json:"amount"`
	Quantity         float64   `json:"quantity"`
	InvestmentDate   time.Time `json:"investment_date"`
	Description      string    `json:"description"`
}

type InvestmentsRequest struct {
	InvestmentTypeID string    `json:"investment_type_id"`
	UserID           string    `json:"user_id"`
	Name             string    `json:"name"`
	Amount           float64   `json:"amount"`
	Quantity         float64   `json:"quantity"`
	InvestmentDate   time.Time `json:"investment_date"`
	Description      string    `json:"description"`
}
