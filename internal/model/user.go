package model

type User struct {
	ID           int    `json:"id" db:"id"`
	Email        string `json:"email" db:"email" binding:"required"`
	PasswordHash []byte `json:"password" db:"passHash" binding:"required"`
}
