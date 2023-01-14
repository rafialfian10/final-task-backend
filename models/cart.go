package models

import "time"

type Cart struct {
	Id            int         `json:"id" gorm:"primary_key:auto_increment"`
	BookId        int         `json:"book_id" gorm:"type: int"`
	Book          Book        `json:"book" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TransactionId int         `json:"transaction_id"`
	Transaction   Transaction `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Total         int         `json:"total"`
	CreateAt      time.Time   `json:"-"`
}
