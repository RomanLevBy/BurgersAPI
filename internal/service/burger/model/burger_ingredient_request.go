package model

type BurgerIngredientRequest struct {
	IngredientId int    `json:"ingredient_id" validate:"required"`
	Instruction  string `json:"instruction" validate:"required"`
}
