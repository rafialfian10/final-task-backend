package transactiondto

import "waysbook/models"

type BookResponseForTransaction struct {
	Id                 int    `json:"id"`
	Title              string `json:"title" gorm:"type: varchar(255)"`
	PublicationDate    string `json:"publicationdate"`
	ISBN               int    `json:"isbn"`
	Pages              int    `json:"pages"`
	Author             string `json:"author"`
	Price              int    `json:"price"`
	IsPromo            bool   `json:"is_promo"`
	Discount           int    `json:"discount"`
	PriceAfterDiscount int    `json:"price_after_discount"`
	Description        string `json:"description" gorm:"type: text"`
	Book               string `json:"book"`
	Thumbnail          string `json:"thumbnail"`
	Quota              int    `json:"quota" form:"quota"`
	OrderQty           int    `json:"orderQty"`
}

type TransactionResponse struct {
	Id         string                       `json:"id"`
	MidtransId string                       `json:"midtrans_id"`
	OrderDate  string                       `json:"order_date"`
	Total      int                          `json:"total"`
	Status     string                       `json:"status"`
	User       models.UserResponse          `json:"user"`
	BookId     int                          `json:"-"`
	Book       []BookResponseForTransaction `json:"book"`
}
