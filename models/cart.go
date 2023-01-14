package models

type Cart struct {
	Id            int         `json:"id" gorm:"primary_key:auto_increment"`
	Qty           int         `json:"qty"`
	BookId        int         `json:"book_id" gorm:"type: int"`
	Book          Book        `json:"book" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserId        int         `json:"user_id"`
	User          User        `json:"user"`
	TransactionId int         `json:"transaction_id"`
	Transaction   Transaction `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Total         int         `json:"subtotal"`
}
