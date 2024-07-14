package model

type BurgerRequest struct {
	CategoryId   int    `json:"category_id" validate:"required"`
	Title        string `json:"title" validate:"required"`
	Instructions string `json:"instructions" validate:"required"`
	Video        string `json:"video,omitempty"`
}
