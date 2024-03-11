package model

type Category struct {
	Name string `json:"name" binding:"required"`
}
