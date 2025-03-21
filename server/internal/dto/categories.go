package dto

type CategoryType string

const (
	Income  CategoryType = "income"
	Expense CategoryType = "expense"
)

type CategoriesResponse struct {
	ID       string       `json:"id"`
	ParentID string       `json:"parent_id"`
	Name     string       `json:"name"`
	Type     CategoryType `json:"type"`
}

type CategoriesRequest struct {
	ParentID string       `json:"parent_id"`
	Name     string       `json:"name"`
	Type     CategoryType `json:"type"`
}
