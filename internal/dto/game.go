package dto

type GameDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name,omitempty"`
}
