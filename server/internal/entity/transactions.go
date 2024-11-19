package entity

type Transactions struct {
	Base
	Amount          int    `gorm:"type:decimal(15)"`
	TransactionType string `gorm:"type:varchar(100);not null"`
	Date            string `gorm:"type:date;not null"`
	Description     string `gorm:"type:text"`
	Category        string `gorm:"type:varchar(100);not null"`
	AttachmentUrl   string `gorm:"type:text"`
	UserID          string `gorm:"type:uuid;not null"`

	User Users `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type TransactionDetail struct {
	ID              string `json:"id"`
	Amount          int    `json:"amount"`
	TransactionType string `json:"transaction_type"`
	Date            string `json:"date"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	UserName        string `json:"user_name"`
	UserEmail       string `json:"user_email"`
}

type TransactionsResponse struct {
	ID              string `json:"id"`
	Amount          int    `json:"amount"`
	TransactionType string `json:"transaction_type"`
	Date            string `json:"date"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	AttachmentUrl   string `json:"attachment_url"`
	UserID          string `json:"user_id"`
}

type TransactionsRequest struct {
	Amount          int    `json:"amount"`
	TransactionType string `json:"transaction_type"`
	Date            string `json:"date"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	AttachmentUrl   string `json:"attachment_url"`
	UserID          string `json:"user_id"`
}
