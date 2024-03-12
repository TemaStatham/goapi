package model

type Product struct {
	ID          int        `json:"id"`
	Name        string     `json:"name" binding:"required"`
	Categoryies []Category `json:"categoryies" binding:"required"`
}
