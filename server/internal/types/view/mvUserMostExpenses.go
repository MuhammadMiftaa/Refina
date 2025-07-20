package view

type MVUserMostExpenses struct {
	UserID             string  `json:"user_id"`
	ParentCategoryName string  `json:"parent_category_name"`
	Total              float64 `json:"total"`
	Rank               int     `json:"rank"`
}
