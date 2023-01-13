package booksdto

type CreateBookRequest struct {
	Title           string `json:"title" form:"title" gorm:"type: varchar(255)"`
	PublicationDate string `json:"publication_date"`
	ISBN            int    `json:"isbn" form:"isbn" gorm:"type: int"`
	Pages           int    `json:"pages" form:"pages" gorm:"type: int"`
	Author          string `json:"author" form:"author" gorm:"type: varchar(255)"`
	Price           int    `json:"price" form:"price" gorm:"type: int"`
	Description     string `json:"description" form:"description" gorm:"type: text"`
	BookAttachment  string `json:"book_attachment" form:"book_attachment"`
	Thumbnail       string `json:"thumbnail" form:"thumbnail"`
}

type UpdateBookRequest struct {
	Title              string `json:"title" gorm:"type: varchar(255)"`
	PublicationDate    string `json:"publication_date"`
	ISBN               int    `json:"isbn"`
	Pages              int    `json:"pages"`
	Author             string `json:"author"`
	Price              int    `json:"price"`
	IsPromo            bool   `json:"is_promo"`
	Discount           int    `json:"discount"`
	PriceAfterDiscount int    `json:"price_after_discount"`
	Description        string `json:"description" gorm:"type: text"`
	BookAttachment     string `json:"book_attachment"`
	Thumbnail          string `json:"thumbnail"`
}
