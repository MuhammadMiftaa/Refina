package dto

type ReportRequest struct {
	UserID   string `json:"user_id"`
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
}
