package converter

import (
	serviceModel "github.com/RomanLevBy/BurgersAPI/internal/service/burger/model"
)

func ToBurgerInfoFromRequest(br serviceModel.BurgerRequest) serviceModel.BurgerInfo {
	return serviceModel.BurgerInfo{
		CategoryId:   br.CategoryId,
		Title:        br.Title,
		Instructions: br.Instructions,
		Video:        br.Video,
	}
}
