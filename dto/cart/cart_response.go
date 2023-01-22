package cartdto

import "waysbook/models"

type CartResponse struct {
	Id            int         `json:"id"`
	BookId        int         `json:"book_id"`
	UserId        int         `json:"user_id"`
	BookTitle     string      `json:"book_title"`
	BookThumbnail string      `json:"book_thumbnail"`
	Author        string      `json:"author"`
	OrderQty      int         `json:"order_qty"`
	Book          models.Book `json:"book"`
}

type DeleteCartResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}
