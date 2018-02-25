package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/kataras/iris/context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris"
	"github.com/pankajojha/go-db-data/backend"

	"bytes"
	"encoding/json"

	"html/template"
)

var db *sql.DB

func main() {
	printData(db)
}

func printData(db *sql.DB) {
	app := iris.New()
	app.RegisterView(iris.HTML("./frontend/templates", ".html"))

	app.Get("/", func(ctx context.Context) {
		ctx.View("home.html")
	})

	app.Post("/query1", func(ctx context.Context) {
		query := ctx.PostValue("query")
		wfDbs := dbutil.GetWfDB("")

		wfBasedDatas := []dbutil.Webdata{}
		for _, col := range wfDbs {
			wfdata := dbutil.QueryDB(col.IpAddress, col.WfSchemaName, query, "Workflow Data", col.Id+" : "+col.WfSchemaName+" : "+col.IpAddress, col.DataSourceName)
			wfBasedDatas = append(wfBasedDatas, wfdata)
		}
		ctx.ViewData("data", wfBasedDatas[0])
		ctx.View("result.html")
	})

	app.Post("/query", func(ctx context.Context) {
		//outputTmpl := template.Must(template.ParseFiles("output.html"))
		outputTmpl := template.Must(template.ParseFiles("./frontend/templates/output.html"))

		query := ctx.PostValue("query")
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

		outputTmpl.Execute(os.Stdout, wfRespond)
	})

	app.Get("/wf", func(ctx context.Context) {
		query := "select * from ps_wf_instance"
		wdata := dbutil.QueryDB("", "", query, "MyTitle", "WF Heading", " Grid title")

		ctx.ViewData("data", wdata)
		//ctx.View("hello.html")
		ctx.View("result.html")
	})

	app.Get("/json", func(ctx context.Context) {
		b := []byte(`{"hello": "123"}`)
		b, _ = prettyprint(b)
		fmt.Printf("%s", b)
	})

	// Start the server using a network address and block.
	app.Run(iris.Addr(":8090"))
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
