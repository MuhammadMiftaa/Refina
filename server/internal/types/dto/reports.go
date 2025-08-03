package dto

import "time"

type ReportRequest struct {
	UserID   string `json:"user_id"`
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
}

type ReportResponse struct {
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	FromDate    time.Time `json:"from_date"`
	ToDate      time.Time `json:"to_date"`
	FileURL     string    `json:"file_url"`
	FileSize    int64     `json:"file_size"`
	GeneratedAt time.Time `json:"generated_at"`
}
