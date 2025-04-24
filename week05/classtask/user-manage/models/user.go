package models

type User struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age"`
}
