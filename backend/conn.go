package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

//DBConfig ...
type DBConfig struct {
	Host              string
	UserName          string
	Password          string
	Port              string
	Db                string
	DbTable           string
	Timeout           string
	MaxOpenConnection int
	MaxIdleConnection int
	LogPath           string
}

var (
	Config DBConfig
)

func init() {
	Config = ReadConfig()
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

//GetDbByIp is used to get connection
func GetDbByIP(ip string, dbName string) *sql.DB {

	if ip == "" {
		ip = Config.Host
	}

	if dbName == "" {
		dbName = Config.Db
	}
	connection := Config.UserName + ":" + Config.Password + "@tcp(" + ip + ":" + Config.Port + ")/" + dbName
	db, err := sql.Open("mysql", connection)

	if Config.MaxOpenConnection != 0 {
		db.SetMaxOpenConns(Config.MaxOpenConnection)
	}

	if Config.MaxIdleConnection != 0 {
		db.SetMaxIdleConns(Config.MaxIdleConnection)
	}

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	return db
}
