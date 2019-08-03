package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/tkanos/gonfig"
)

type DatabaseConfig struct {
	Host        string    `json:"host"`
    Port        int       `json:"port"`
    User        string    `json:"user"`
    Password    string    `json:"password"`
    Database    string    `json:"database"`
}

func InitDB() (*sql.DB) {
	databaseConfig := DatabaseConfig{}
	err := gonfig.GetConf("config/config.production.json", &databaseConfig)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		databaseConfig.Host, databaseConfig.Port, databaseConfig.User, 
		databaseConfig.Password, databaseConfig.Database)

	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	return database
}
