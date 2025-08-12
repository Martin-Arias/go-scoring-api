package dto

type CreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type GameResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
