package view

type Attachment struct {
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	Image         string `json:"image"`
	Format        string `json:"format"`
	Size          int64  `json:"size"`
}

type ViewUserTransactions struct {
	ID              string       `json:"id"`
	UserID          string       `json:"user_id"`
	WalletID        string       `json:"wallet_id"`
	WalletNumber    string       `json:"wallet_number"`
	WalletType      string       `json:"wallet_type"`
	WalletTypeName  string       `json:"wallet_type_name"`
	WalletBalance   float64      `json:"wallet_balance"`
	CategoryID      string       `json:"category_id"`
	CategoryName    string       `json:"category_name"`
	CategoryType    string       `json:"category_type"`
	Amount          float64      `json:"amount"`
	TransactionDate string       `json:"transaction_date"`
	Description     string       `json:"description"`
	Attachments     []Attachment `json:"attachments" gorm:"-"`
}
