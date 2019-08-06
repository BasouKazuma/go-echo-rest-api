package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host		string		`mapstructure:"host"`
	Port		int			`mapstructure:"port"`
	User		string		`mapstructure:"user"`
	Password	string		`mapstructure:"password"`
	Database	string		`mapstructure:"database"`
}

func InitDB() (*sql.DB) {
	// Load Config
	databaseConfig := DatabaseConfig{}
	viperDBConfig := viper.Sub("services.database")
	err := viperDBConfig.Unmarshal(&databaseConfig)
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
