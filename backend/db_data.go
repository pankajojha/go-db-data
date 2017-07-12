package dbutil

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//GetAllRows ....
func GetAllRows(db *sql.DB, queryToRun string) *sql.Rows {
	defer db.Close()
	// Execute the query
	rows, err := db.Query(queryToRun)
	checkErr(err)
	return rows
}

//GetData Gets data from the server
func GetData(db *sql.DB, queryToRun string) {

	//db := GetDB(connectionString)
	defer db.Close()
	// Execute the query
	rows, err := db.Query(queryToRun)
	checkErr(err)

	// Get column names
	columns, err := rows.Columns()
	checkErr(err)

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		checkErr(err)
		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
