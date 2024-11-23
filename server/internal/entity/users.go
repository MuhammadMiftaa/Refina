package entity

import "database/sql"

type Users struct {
	Base
	Name           string       `gorm:"type:varchar(100);not null"`
	Email          string       `gorm:"type:varchar(100);unique;not null"`
	Password       string       `gorm:"type:varchar(100);not null"`
	EmailVerfiedAt sql.NullTime
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
