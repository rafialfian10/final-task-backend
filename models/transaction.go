package models

import "time"

type Transaction struct {
	ID        int       `json:"id" gorm:"primary_key:auto_increment"`
	Total     int       `json:"total"`
	UserID    int       `json:"user_id"` // who doing the transaction
	User      User      `json:"user"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type BookTransaction struct {
	ID int `json:"id"`
}

type UserTransaction struct {
	ID int `json:"id"`
}

func (UserTransaction) TableName() string {
	return "transactions"
}
