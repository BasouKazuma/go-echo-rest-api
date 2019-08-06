package main

import (
	"fmt"
	"net/http"
	
	"go-echo-rest-api/db"
	"go-echo-rest-api/handler"
	
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

func main() {
	// Load Config
	viper.SetConfigName("viper.production") // name of config file (without extension)
	viper.AddConfigPath("config/") // path to look for the config file in
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Initialize Echo
	e := echo.New()

	//Initialize Database
	database := db.InitDB()
	// defer database.Close()

	// Initialize handler
	h := &handler.Handler{DB: database}

	// Routes
	e.GET("/", getHome)
	e.POST("/user/create", h.CreateUser)
	e.GET("/user/:user_id", h.GetUser)
	e.POST("/user/:user_id/file/create", h.CreateFile)
	e.GET("/user/:user_id/file/list", h.GetFileList)
	e.GET("/user/:user_id/file/:file_hash", h.GetFile)

	e.Logger.Fatal(e.Start(":1323"))
}

func getHome(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}
