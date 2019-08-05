package handler

import (
	"go-echo-rest-api/model"
	"net/http"

	"github.com/labstack/echo"
)

type CreateUserBody struct {
	Email	string	`json:"email" db:"email"`
}

func (h *Handler) CreateUser(c echo.Context) (err error) {
	// Bind
	userBody := new(CreateUserBody)
	if err = c.Bind(userBody); err != nil {
		return
	}

	// Check Connection
	database := h.DB
	err = database.Ping()
	if err != nil {
		panic(err)
	}

	// Check email
	var count int
	checkUserQuery := `SELECT COUNT(*) FROM users WHERE email = $1`
	err = database.QueryRow(checkUserQuery, userBody.Email).Scan(&count)
	if count > 0 {
		message := model.ErrorMessage{Message: "Email already exists."}
		return c.JSON(http.StatusConflict, message)
	}
	
	// Add User
	user := model.User{Email: userBody.Email}
	addUserQuery := `INSERT INTO users ( email ) 
 		VALUES( $1 )
 		RETURNING id`
	err = database.QueryRow(addUserQuery, userBody.Email).Scan(&user.Id)
	if err != nil {
		message := model.ErrorMessage{Message: "An error occurred."}
		return c.JSON(http.StatusBadRequest, message)
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) GetUser(c echo.Context) (err error) {
	userId := c.Param("user_id")

	// Check Connection
	database := h.DB
	err = database.Ping()
	if err != nil {
		panic(err)
	}

	// Get new User
	user := model.User{}
	findUserQuery := `SELECT * FROM users WHERE id = $1`
	err = database.QueryRow(findUserQuery, userId).Scan(&user.Id, &user.Email)
	if err != nil {
		message := model.ErrorMessage{Message: "Account not found."}
		return c.JSON(http.StatusNotFound, message)
	}
	return c.JSON(http.StatusOK, user)
}
