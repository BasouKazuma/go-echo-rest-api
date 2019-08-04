package model

type User struct {
	Id       int64     `json:"id" db:"id"`
	Email    string    `json:"email" db:"email"`
}

type ErrorMessage struct {
	Message    string    `json:"message"`
}
