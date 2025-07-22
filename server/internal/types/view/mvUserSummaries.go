package view

type MVUserSummaries struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	UserIncome   float64 `json:"user_income"`
	UserExpense  float64 `json:"user_expense"`
	UserProfit   float64 `json:"user_profit"`
	TotalBalance float64 `json:"total_balance"`
}
