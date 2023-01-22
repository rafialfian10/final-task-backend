package cartdto

type CreateCartRequest struct {
	BookId int `json:"book_id"`
}

type UpdateCartRequest struct {
	Event string `json:"event"`
	Qty   int    `json:"qty"`
}
