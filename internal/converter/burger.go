package converter

import (
	serviceModel "github.com/RomanLevBy/BurgersAPI/internal/service/burger/model"
)

func ToBurgerInfoFromRequest(br serviceModel.BurgerRequest) serviceModel.BurgerInfo {
	ingredients := make([]serviceModel.BurgerIngredientInfo, len(br.Ingredients))

	for i, ing := range br.Ingredients {
		ingredient := serviceModel.BurgerIngredientInfo{
			IngredientId: ing.IngredientId,
			Instruction:  ing.Instruction,
		}

		ingredients[i] = ingredient
	}

	return serviceModel.BurgerInfo{
		CategoryId:   br.CategoryId,
		Title:        br.Title,
		Instructions: br.Instructions,
		Video:        br.Video,
		Ingredients:  ingredients,
	}
}
