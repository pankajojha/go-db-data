package dbutil

import (
	"database/sql"
	"fmt"
)

type Data struct {
	Value string
}

type Webdata struct {
	Title string
	Heading string
	GridTitle string
	ColumnHeading []string
	RowData []Data
	NumOfRows *sql.Rows
}

func (webData *Webdata) AddColumnHeading(data string) []string {
	webData.ColumnHeading = append(webData.ColumnHeading, data)
	return webData.ColumnHeading
}

func (webData *Webdata) AddRowData(data Data) []Data {
	webData.RowData = append(webData.RowData, data)
	return webData.RowData
}


func QueryDB(ip string, schema string, query string, title string, heading string, gridTitle string) Webdata {

	db := GetDbByIp(ip, schema)

	wdata := Webdata{Title: title, Heading: heading, GridTitle: gridTitle}
	rows :=GetData(db, query)
	columns, err := rows.Columns()
	CheckErr(err)


	// Make a slice for the values
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}
	//fill table headings, the table returns 9 columns so I just hard coded it
	for j:=0;j<len(columns);j++ {
		wdata.AddColumnHeading(columns[j])
	}

	wdata.NumOfRows = rows

	// Fetch rows
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}
		// Print data
		vData := Data{}

		var rowData string
		for i, col := range values {
			switch col.(type) {
			case nil:
				rowData +="NULL "
			case []byte:
				rowData +=  string(col.([]byte))
			default:
				rowData +=  fmt.Sprint(col)
			}
			rowData += " | "
			vData.Value = rowData

			fmt.Println(columns[i], ": ", rowData)
		}
		wdata.AddRowData(vData)

	}
	return wdata
}

//GetData Gets data from the server
func GetData(db *sql.DB, queryToRun string) *sql.Rows {
	//db := GetDB(connectionString)
	defer db.Close()
	// Execute the query
	rows, err := db.Query(queryToRun)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return rows
}

//CheckErr ....
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
