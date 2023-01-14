package models

type Transaction struct {
	Id         int              `json:"id" gorm:"primary_key:auto_increment"`
	CounterQty int              `json:"qty" form:"qty" gorm:"type: int"`
	UserId     int              `json:"user_id"`
	User       UserResponse     `json:"user"`
	Cart       []Cart           `json:"cart"`
	BookId     int              `json:"book_id"`
	Book       BookCartResponse `json:"book"`
	Attachment string           `json:"attachment"`
	Total      int              `json:"total"`
	Status     string           `json:"status"`
}
