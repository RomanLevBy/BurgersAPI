package model

type Burger struct {
	ID           int    `json:"id"`
	Category     string `json:"category"`
	Handle       string `json:"handler"`
	Title        string `json:"title"`
	Instructions string `json:"instructions"`
	Video        string `json:"video"`
	DataModified string `json:"data_modified"`
}
