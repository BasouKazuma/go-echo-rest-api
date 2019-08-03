package handler

import (
	"go-echo-rest-api/model"
	"net/http"

	"github.com/labstack/echo"
)

func (h *Handler) CreateUser(c echo.Context) (err error) {
	// // Bind
	// user := new(model.User)
	// if err = c.Bind(user); err != nil {
	// 	return
	// }

	// Data
	database := h.DB
	email := "test@test.com"
	var counter int
	database.QueryRow("SELECT count(*) FROM users WHERE email = $1", email).Scan(&counter)
	if counter > 0 {
		message := model.ErrorMessage{Message: "Email already exists."}
		return c.JSON(http.StatusConflict, message)
	}
	user := model.User{}
	err = database.QueryRow("SELECT * FROM users WHERE email = $1", email).Scan(&user)
	if err != nil {
		message := model.ErrorMessage{Message: "An error occurred."}
		return c.JSON(http.StatusConflict, message)
	}
	return c.JSON(http.StatusCreated, user)
}
