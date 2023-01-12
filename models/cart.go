package models

type Cart struct {
	ID       int              `json:"id" gorm:"primary_key:auto_increment"`
	QTY      int              `json:"qty"`
	Subtotal int              `json:"subtotal"`
	SellerID int              `json:"seller_id"`
	BookID   int              `json:"book_id"  gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Book     BookCartResponse `json:"book" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID   int              `json:"user_id"`
	User     User             `json:"user"`
}

type UserCart struct {
	ID int `json:"id"`
}

func (UserCart) TableName() string {
	return "carts"
}
