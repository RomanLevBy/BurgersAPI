package model

type Category struct {
	ID     int    `json:"id"`
	Handle string `json:"handler"`
	Title  string `json:"title"`
}
