package dto

type InvestmentTypesResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type InvestmentTypesRequest struct {
	Name string `json:"name"`
	Unit string `json:"unit"`
}
