package bookdto

type BookResponse struct {
	Id                 int    `json:"id"`
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
	Quota              int    `json:"quota" form:"quota"`
}
