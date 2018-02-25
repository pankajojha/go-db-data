package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"github.com/pankajojha/go-db-data/backend"
)

var db *sql.DB

func main() {

	homeTmpl := template.Must(template.ParseFiles("./frontend/templates/home.html"))
	outputTmpl := template.Must(template.ParseFiles("./frontend/templates/output.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeTmpl.Execute(w, "")
	})

	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		query := r.FormValue("query")
		fmt.Println(" query ", query)

		wfDbs := dbutil.GetWfDB("")
		noOfWorkflows := len(wfDbs)

		var wg sync.WaitGroup
		wfRespond := make(chan dbutil.Webdata, noOfWorkflows)
		wg.Add(noOfWorkflows)

		for _, col := range wfDbs {
			go dbutil.GetDbDataOnChan(wfRespond, &wg, query, "Workflow Data", col)
		}

		wg.Wait()
		close(wfRespond)

		outputTmpl.Execute(w, wfRespond)
	})

	http.ListenAndServe(":8090", nil)
}
