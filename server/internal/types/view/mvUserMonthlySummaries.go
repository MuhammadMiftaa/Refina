package view

type MVUserMonthlySummaries struct {
	UserID       string  `json:"user_id"`
	Month        string  `json:"month"`
	MonthName    string  `json:"month_name"`
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
}
