package model

type User struct {
	Id       int64     `json:"id"`
	Email    string    `json:"email"`
}

type ErrorMessage struct {
	Message    string    `json:"message"`
}
