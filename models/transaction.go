package models

import "time"

type Transaction struct {
	Id         string         `json:"id" gorm:"primary_key:auto_increment"`
	MidtransId string         `json:"midtrans_id" gorm:"type: varchar(255)"`
	OrderDate  time.Time      `json:"order_date"`
	UserId     int            `json:"user_id"`
	User       UserResponse   `json:"user"`
	Cart       []CartResponse `json:"cart" gorm:"foreignKey:TransactionId"`
	Total      int            `json:"total"`
	Status     string         `json:"status"`
}

type TransactionResponse struct {
	Id         string       `json:"id" gorm:"type: varchar(255);PRIMARY_KEY"`
	MidtransId string       `json:"midtrans_id" gorm:"type: varchar(255)"`
	OrderDate  time.Time    `json:"order_date"`
	UserId     int          `json:"user_id"`
	User       UserResponse `json:"users"`
	Total      int          `json:"total"`
	Status     string       `json:"status"`
}

func (TransactionResponse) TableName() string {
	return "transactions"
}
