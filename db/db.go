package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host		string		`json:"host"`
	Port		int			`json:"port"`
	User		string		`json:"user"`
	Password	string		`json:"password"`
	Database	string		`json:"database"`
}

func InitDB() (*sql.DB) {
	// Load Config
	databaseConfig := DatabaseConfig{}
	viper.SetConfigName("database") // name of config file (without extension)
	viper.AddConfigPath("config/") // path to look for the config file in
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = viper.Unmarshal(&databaseConfig)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	// Setup Connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		databaseConfig.Host,
		databaseConfig.Port,
		databaseConfig.User, 
		databaseConfig.Password,
		databaseConfig.Database)
	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Check Connection
	err = database.Ping()
	if err != nil {
	 panic(err)
	}
	return database
}
