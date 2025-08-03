package entity

import "time"

type Reports struct {
	Base

	UserID        string    `gorm:"type:uuid;not null" json:"user_id"`
	FromDate      time.Time `gorm:"type:timestamptz;not null" json:"from_date"`
	ToDate        time.Time `gorm:"type:timestamptz;not null" json:"to_date"`
	RequestAt     time.Time `gorm:"type:timestamptz;not null" json:"request_at"`
	NextRequestAt time.Time `gorm:"type:timestamptz;not null" json:"next_request_at"`
	Status        string    `gorm:"type:report_status;not null" json:"status"`
	FileURL       *string    `gorm:"type:text" json:"file_url"`
	FileSize      *int64     `gorm:"type:bigint" json:"file_size"`
	GeneratedAt   *time.Time `gorm:"type:timestamptz" json:"generated_at"`
}
