package model

type User struct {
	ID       int    `json:"id" db:"id"`
	Email    string `json:"email" db:"email" binding:"required"`
	PassHash []byte `json:"password" db:"passhash" binding:"required"`
}
