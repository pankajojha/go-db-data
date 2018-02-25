package dbutil

import (
	"database/sql"
	"fmt"
	"sync"
)

type Data struct {
	Value string
}

type Datas struct {
	Value []string
}

type Webdata struct {
	Title         string
	Heading       string
	GridTitle     string
	ColumnHeading []string
	RowData       []Data
	RowDatas      []Datas
	ValueData     []interface{}
	NumOfRows     *sql.Rows
}

type WfDb struct {
	Id             string
	WfSchemaName   string
	IpAddress      string
	DataSourceName string
	AccountName    string
}

func (webData *Webdata) AddColumnHeading(data string) []string {
	webData.ColumnHeading = append(webData.ColumnHeading, data)
	return webData.ColumnHeading
}

func (webData *Webdata) AddRowData(data Data) []Data {
	webData.RowData = append(webData.RowData, data)
	return webData.RowData
}

func (webData *Webdata) AddRowsData(data Datas) []Datas {
	webData.RowDatas = append(webData.RowDatas, data)
	return webData.RowDatas
}

func QueryDBChannel(ip string, schema string, query string, title string, heading string, gridTitle string, c chan Webdata) {
	webData := QueryDB(ip, schema, query, title, heading, gridTitle)
	c <- webData
}

func GetDbDataOnChan(wfRespond chan<- Webdata, wg *sync.WaitGroup, query string, ns string, col WfDb) {
	defer wg.Done()
	wfRespond <- QueryDB(col.IpAddress, col.WfSchemaName, query, ns, col.Id+" : "+col.WfSchemaName+" : "+col.IpAddress, col.DataSourceName)
}

func QueryDB(ip string, schema string, query string, title string, heading string, gridTitle string) Webdata {

	db := GetDbByIP(ip, schema)

	wdata := Webdata{Title: title, Heading: heading, GridTitle: gridTitle}
	rows, err := GetData(db, query)

	if err != nil {
		errorData := Data{" Error in query : " + query}
		wdata.AddRowData(errorData)
		return wdata
	}

	columns, err := rows.Columns()
	CheckErr(err)

	// Make a slice for the values
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}
	//fill table headings, the table returns 9 columns so I just hard coded it
	for j := 0; j < len(columns); j++ {
		wdata.AddColumnHeading(columns[j])
	}

	wdata.NumOfRows = rows

	// Fetch rows
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		CheckErr(err)
		// Print data
		vData := Data{}

		//	datas :=  make([]Datas{}, len(columns))
		datas := make([]string, len(columns))
		vDatas := Datas{}

		var rowData string
		for i, col := range values {
			switch col.(type) {
			case nil:
				rowData = "NULL "
			case []byte:
				rowData = string(col.([]byte))
			default:
				rowData = fmt.Sprint(col)
			}
			//rowData += " | "
			vData.Value = rowData
			datas[i] = rowData
			//fmt.Println(columns[i], ": ", rowData)
		}
		vDatas.Value = datas
		//wdata.AddRowData(vData)
		wdata.AddRowsData(vDatas)
		wdata.ValueData = values

	}
	return wdata
}

//GetData Gets data from the server
func GetData(db *sql.DB, queryToRun string) (*sql.Rows, error) {
	//db := GetDB(connectionString)
	defer db.Close()
	// Execute the query
	rows, err := db.Query(queryToRun)
	if err != nil {
		return rows, fmt.Errorf("sql: can not query with query  "+queryToRun+"  %v", err)
	}

	return rows, err
}

//CheckErr ....
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetWfDB(ip string) []WfDb {

	config := ReadConfig()
	db := GetDbByIP(ip, "epenops")

	query := "select * from " + config.DbTable

	rows, err := GetData(db, query)
	columns, err := rows.Columns()
	CheckErr(err)

	// Make a slice for the values
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	wfDbs := []WfDb{}

	// Fetch rows
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		CheckErr(err)
		// Print data
		wfDb := WfDb{}

		var rowData string
		for i, col := range values {
			switch col.(type) {
			case nil:
				rowData = "NULL "
			case []byte:
				rowData = string(col.([]byte))
			default:
				rowData = fmt.Sprint(col)
			}
			if columns[i] == "ps_wf_instance_id" {
				wfDb.Id = rowData
			} else if columns[i] == "ip_address" {
				wfDb.IpAddress = rowData
			} else if columns[i] == "schema_name" {
				wfDb.WfSchemaName = rowData
			} else if columns[i] == "account_name" {
				wfDb.AccountName = rowData
			} else if columns[i] == "datasource_name" {
				wfDb.DataSourceName = rowData
			}
		}
		wfDbs = append(wfDbs, wfDb)
	}
	return wfDbs
}
