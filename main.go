package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pankajojha/go-db-data/backend"
)

func main() {
	config := dbutil.ReadConfig()
	connectionString := config.UserName + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.Db
	query := "select * from items"
	fmt.Println(connectionString)
	fmt.Println(query)

	db := dbutil.GetDB(connectionString)

	dbutil.GetData(db, query)
}
