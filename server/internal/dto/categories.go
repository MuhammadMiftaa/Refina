package dto

type CategoryType string

const (
	Income  CategoryType = "income"
	Expense CategoryType = "expense"
)

type Category struct {
	ID   string       `json:"id"`
	Name string       `json:"name"`
}

type CategoriesResponse struct {
	GroupName string     `json:"group_name"`
	Category  []Category `json:"category"`
	Type CategoryType `json:"type"`
}

type CategoriesRequest struct {
	ParentID string       `json:"parent_id"`
	Name     string       `json:"name"`
	Type     CategoryType `json:"type"`
}
