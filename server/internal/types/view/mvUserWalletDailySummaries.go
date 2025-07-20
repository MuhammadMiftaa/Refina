package view

type MVUserWalletDailySummaries struct {
	UserID   string  `json:"user_id"`
	Date     string  `json:"date"`
	Physical float64 `json:"physical"`
	EWallet  float64 `json:"e-wallet"`
	Bank     float64 `json:"bank"`
	Others   float64 `json:"others"`
}
