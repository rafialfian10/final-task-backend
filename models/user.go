package models

type User struct {
	Id        int    `json:"id" gorm:"primary_key:auto_increment"`
	Name      string `json:"name" gorm:"type: varchar(255)"`
	Email     string `json:"email" gorm:"type: varchar(255)"`
	Password  string `json:"password" gorm:"type: varchar(255)"`
	Gender    string `json:"gender" gorm:"type: varchar(255)"`
	Phone     string `json:"phone" gorm:"type: varchar(255)"`
	Thumbnail string `json:"thumbnail"`
	Address   string `json:"address" gorm:"type: text"`
	Role      string `json:"role" gorm:"type: varchar(255)"`
}

type UserResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Gender    string `json:"gender"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Thumbnail string `json:"thumbnail"`
}

func (UserResponse) TableName() string {
	return "users"
}
