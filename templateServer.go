package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/pankajojha/go-db-data/backend"
)

var templates *template.Template

//getRows
func getRows(query string) *sql.Rows {
	config := dbutil.ReadConfig()
	connectionString := config.UserName + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.Db
	db := dbutil.GetDB(connectionString)
	rows := dbutil.GetAllRows(db, query)
	return rows
}

func handler1(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("t1.html", "t2.html")
	t.Execute(w, "Asit")
}

func handler2(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("t1.html", "t2.html")
	t.ExecuteTemplate(w, "t2.html", "Golang")
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	query := "select * from items"

	rows := getRows(query)
	columns, _ := rows.Columns()

	templates := parseTemplate()
	s1 := templates.Lookup("header.tmpl")
	s1.ExecuteTemplate(os.Stdout, "header", columns)

	s2 := templates.Lookup("content.tmpl")
	s2.ExecuteTemplate(os.Stdout, "content", nil)

	s3 := templates.Lookup("footer.tmpl")
	s3.ExecuteTemplate(os.Stdout, "footer", nil)
	fmt.Println()
	s3.Execute(os.Stdout, nil)

	http.HandleFunc("/t1", handler1)
	http.HandleFunc("/t2", handler2)
	server.ListenAndServe()
}

//parseTemplate ...
func parseTemplate() *template.Template {
	var allFiles []string
	files, err := ioutil.ReadDir("./templates")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".tmpl") {
			allFiles = append(allFiles, "./templates/"+filename)
		}
	}

	templates, err = template.ParseFiles(allFiles...)

	if err != nil {
		print(err)
	}

	return templates
}
