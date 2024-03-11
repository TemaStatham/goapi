package model

type Product struct {
	Name        string     `json:"name" binding:"required"`
	Categoryies []Category `json:"categoryies" binding:"required"`
}
