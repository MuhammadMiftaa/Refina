package entity

import "database/sql"

type Role string

const (
	Admin Role = "admin"
	User  Role = "user"
)

type Users struct {
	Base
	Name           string       `gorm:"type:varchar(100);not null"`
	Email          string       `gorm:"type:varchar(100);unique;not null"`
	Password       string       `gorm:"type:varchar(100);not null"`
	Role           string       `gorm:"type:varchar(100);not null;default:'user'"`
	EmailVerfiedAt sql.NullTime `gorm:"type:timestamp"`
}

type UserWallet struct {
	ID            string
	UserID        string
	Name          string
	Email         string
	WalletNumber  string
	WalletBalance float64
	WalletName    string
	WalletType    string
}

type UserInvestment struct {
	ID                 string
	UserID             string
	Name               string
	Email              string
	InvestmentType     string
	InvestmentName     string
	InvestmentAmount   float64
	InvestmentQuantity float64
	InvestmentUnit     string
	InvestmentDate     string
}

type UserTransactions struct {
	UserID string
	Name   string
	Email  string

	WalletID      string
	WalletNumber  string
	WalletBalance float64
	WalletName    string
	WalletType    string

	TransactionID   string
	CategoryName    string
	CategoryType    string
	Amount          float64
	TransactionDate string
	Description     string

	Image string
}

type UsersResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UsersRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
