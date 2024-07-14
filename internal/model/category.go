package model

type Category struct {
	ID      int    `json:"id"`
	Handler string `json:"handler"`
	Title   string `json:"title"`
}
