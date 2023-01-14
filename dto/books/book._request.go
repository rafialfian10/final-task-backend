package bookdto

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
	Quota           int    `json:"quota" form:"quota"`
}

type UpdateBookRequest struct {
	Id       string `json:"id"`
	Discount string `json:"discount"`
}
