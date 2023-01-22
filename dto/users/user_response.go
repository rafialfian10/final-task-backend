package userdto

type UserResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Thumbnail string `json:"thumbnail"`
	Role      string `json:"role" form:"role"`
}
