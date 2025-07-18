package dto

type RegisterResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
