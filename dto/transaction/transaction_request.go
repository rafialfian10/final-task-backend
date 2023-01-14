package transactiondto

type CreateTransactionRequest struct {
	CounterQty int `json:"qty" form:"qty"`
	UserId     int `json:"user_id"`
	BookId     int `json:"book_id"`
	// Attachment string `json:"attachment"`
	Total  int    `json:"total"`
	Status string `json:"status"`
}
