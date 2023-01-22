package models

type Cart struct {
	Id            int         `json:"id" gorm:"primary_key:auto_increment"`
	UserId        int         `json:"user_id" gorm:"type: int"`
	User          User        `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BookId        int         `json:"book_id" gorm:"type: int"`
	Book          Book        `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TransactionId int         `json:"transaction_id"`
	Transaction   Transaction `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OrderQty      int         `json:"order_qty" gorm:"type: int"`
}

type CartResponse struct {
	Id            int    `json:"-"`
	TransactionId string `json:"-" gorm:"type: varchar(255)"`
	BookId        int    `json:"-"`
	Book          BookResponse
	OrderQty      int `json:"orderQty" gorm:"type: int"`
}

func (CartResponse) TableName() string {
	return "carts"
}

// constraint digunakan untuk membatasi jenis data yang dapat masuk ke tabel
// fungsi onUpdate Ondelete cascade adalah sebuah fitur yang diberikan untuk sebuah tabel yang berelasi  yang memungkinkan untuk menghapus / update data pada tabel anak apabila data pada tabel parent terhapus /update
