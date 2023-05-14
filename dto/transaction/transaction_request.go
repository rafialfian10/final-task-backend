package transactiondto

type BookRequestForTransaction struct {
	Id       int `json:"id"`
	BookId   int `json:"book_id"`
	OrderQty int `json:"order_qty"`
}

type CreateTransactionRequest struct {
	Total  int                         `json:"total"`
	UserId int                         `json:"user_id"`
	Books  []BookRequestForTransaction `json:"books"`
	// Image  string                      `json:"image" form:"image" gorm:"type: varchar(255)"`
}
type UpdateTransactionRequest struct {
	Total  int                         `json:"total"`
	UserId int                         `json:"user_id"`
	Books  []BookRequestForTransaction `json:"books"`
	Status string                      `json:"status"`
	// Image  string                      `json:"image" form:"image" gorm:"type: varchar(255)"`
}

type UpdateTransactionByAdminRequest struct {
	Status string `json:"status" form:"status"`
}
