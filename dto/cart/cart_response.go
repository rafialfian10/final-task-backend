package cartdto

type CartResponse struct {
	Id            int    `json:"id"`
	Qty           int    `json:"qty"`
	Subtotal      int    `json:"subtotal"`
	BookId        int    `json:"book_id"`
	UserId        int    `json:"user_id"`
	BookTitle     string `json:"book_title"`
	BookThumbnail string `json:"book_thumbnail"`
	Author        string `json:"author"`
}

type DeleteCartResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}
