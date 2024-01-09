package userdto

type UserResponse struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Gender  string `json:"gender"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Photo   string `json:"photo"`
	Role    string `json:"role"`
}
