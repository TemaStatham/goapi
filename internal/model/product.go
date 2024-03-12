package model

type Product struct {
	ID          int        `json:"id" db:"id"`
	Name        string     `json:"name" db:"name" binding:"required"`
	Categoryies []Category `json:"categoryies" binding:"required"`
}
