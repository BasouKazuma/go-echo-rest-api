package main

import (
	"go-echo-rest-api/db"
	"go-echo-rest-api/handler"
	"net/http"
	
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	database := db.InitDB()
	// defer database.Close()

	// Initialize handler
	h := &handler.Handler{DB: database}

	// Routes
	e.GET("/", getHome)
	e.POST("/user/create", h.CreateUser)
	e.GET("/user/:userId", h.GetUser)

	e.Logger.Fatal(e.Start(":1323"))
}

func getHome(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}
