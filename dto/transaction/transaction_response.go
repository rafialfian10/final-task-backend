package transactiondto

import "waysbook/models"

type TransactionResponse struct {
	Id         int                     `json:"id"`
	CounterQty int                     `json:"qty" form:"qty"`
	UserId     int                     `json:"-"`
	User       models.UserResponse     `json:"user"`
	BookId     int                     `json:"-"`
	Book       models.BookCartResponse `json:"book_purchased"`
	Total      int                     `json:"total"`
	Status     string                  `json:"status"`
}
