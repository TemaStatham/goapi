package model

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email" binding:"required"`
	PasswordHash []byte `json:"password" binding:"required"`
}
