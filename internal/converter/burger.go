package converter

import (
	"github.com/RomanLevBy/BurgersAPI/internal/model"
)

func ToBurgerInfoFromRequest(br model.BurgerRequest) model.BurgerInfo {
	return model.BurgerInfo{
		CategoryId:   br.CategoryId,
		Title:        br.Title,
		Instructions: br.Instructions,
		Video:        br.Video,
	}
}
