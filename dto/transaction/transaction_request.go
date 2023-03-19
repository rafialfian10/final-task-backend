package transactiondto

type BookRequestForTransaction struct {
	Id       int `json:"id"`
	BookId   int `json:"book_id"`
	OrderQty int `json:"order_qty"`
}

type CreateTransactionRequest struct {
	Total  int                         `json:"total" validate:"required"`
	UserId int                         `json:"user_id" validate:"required"`
	Books  []BookRequestForTransaction `json:"books" validate:"required"`
}
type UpdateTransactionRequest struct {
	Total  int                         `json:"total"`
	UserId int                         `json:"user_id"`
	Books  []BookRequestForTransaction `json:"books"`
	Status string                      `json:"status"`
}

type UpdateTransactionByAdminRequest struct {
	Status string `json:"status" form:"status"`
}
