package model

import "time"

type BurgerInfo struct {
	CategoryId   int
	Title        string
	Handle       string
	Instructions string
	Video        string
	DataModified time.Time
}
