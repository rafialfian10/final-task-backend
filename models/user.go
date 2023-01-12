package models

type User struct {
	Id       int    `json:"id" gorm:"primary_key:auto_increment"`
	Name     string `json:"name" gorm:"type: varchar(255)"`
	Email    string `json:"email" gorm:"type: varchar(255)"`
	Password string `json:"password" gorm:"type: varchar(255)"`
	Gender   string `json:"gender" gorm:"type: varchar(255)"`
	Phone    string `json:"phone" gorm:"type: varchar(255)"`
	Image    string `json:"image"`
	Address  string `json:"address" gorm:"type: text"`
	Role     string `json:"role" gorm:"type: varchar(255)"`
	Books    []Book `json:"books"`
	// Carts        []Cart        `json:"carts"`
	// Transactions []Transaction `json:"transactions"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Image    string `json:"image"`
}

func (UserResponse) TableName() string {
	return "users"
}
