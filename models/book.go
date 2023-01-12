package models

import "time"

type Book struct {
	Id                 int       `json:"id" gorm:"primary_key:auto_increment"`
	Title              string    `json:"title" gorm:"type: varchar(255)"`
	PublicationDate    time.Time `json:"publication_date"`
	ISBN               int       `json:"isbn"`
	Pages              int       `json:"pages"`
	Author             string    `json:"author"`
	Price              int       `json:"price"`
	IsPromo            bool      `json:"is_promo"`
	Discount           int       `json:"discount"`
	PriceAfterDiscount int       `json:"price_after_discount"`
	Description        string    `json:"description" gorm:"type: text"`
	BookAttachment     string    `json:"book_attachment"`
	Thumbnail          string    `json:"thumbnail"`
	UserId             int       `json:"user_id"`
	User               User      `json:"user"`
}

type BookCartResponse struct {
	Id                 int       `json:"id"`
	Title              string    `json:"title" gorm:"type: varchar(255)"`
	PublicationDate    time.Time `json:"publication_date"`
	ISBN               int       `json:"isbn"`
	Pages              int       `json:"pages"`
	Author             string    `json:"author"`
	Price              int       `json:"price"`
	IsPromo            bool      `json:"is_promo"`
	Discount           int       `json:"discount"`
	PriceAfterDiscount int       `json:"price_after_discount"`
	Description        string    `json:"description" gorm:"type: text"`
	BookAttachment     string    `json:"book_attachment"`
	Thumbnail          string    `json:"thumbnail"`
}

func (BookCartResponse) TableName() string {
	return "books"
}
