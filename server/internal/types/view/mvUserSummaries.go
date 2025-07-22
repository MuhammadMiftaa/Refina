package view

type MVUserSummaries struct {
	UserID                      string  `json:"user_id"`
	Name                        string  `json:"name"`
	IncomeNow                   float64 `json:"income_now"`
	ExpenseNow                  float64 `json:"expense_now"`
	ProfitNow                   float64 `json:"profit_now"`
	BalanceNow                  float64 `json:"balance_now"`
	UserIncomeGrowthPercentage  float64 `json:"user_income_growth_percentage"`
	UserExpenseGrowthPercentage float64 `json:"user_expense_growth_percentage"`
	UserProfitGrowthPercentage  float64 `json:"user_profit_growth_percentage"`
	UserBalanceGrowthPercentage float64 `json:"user_balance_growth_percentage"`
}
