package dbutil

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

//DBConfig ...
type DBConfig struct {
	Host     string
	UserName string
	Password string
	Port     string
	Db       string
}

//ReadConfig ...
func ReadConfig() DBConfig {
	file, _ := os.Open("./conf.json")
	decoder := json.NewDecoder(file)
	configuration := DBConfig{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	return configuration
}

//GetDB is used to get connection
func GetDB(connection string) *sql.DB {
	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	return db
}