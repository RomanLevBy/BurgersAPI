package model

type FetchParam struct {
	Title          string
	TitlePath      string
	IngredientPath string
	Limit          uint64
	CursorID       uint64
}
