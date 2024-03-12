package model

type Category struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name" binding:"required" `
}
