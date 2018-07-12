package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pankajojha/go-db-data/backend"
	"github.com/pankajojha/go-db-data/backend/logger"
)

var db *sql.DB

const (
	VERSION = "0.01"
)

func main() {

	// public views
	http.HandleFunc("/login", LoginFunc)
	// private views
	http.HandleFunc("/query", QueryFunc)
	http.HandleFunc("/wf", QueryWfFunc)
	//http.HandleFunc("/query", PostOnly(BasicAuth(QueryFunc)))
	//http.HandleFunc("/wf", PostOnly(BasicAuth(QueryWfFunc)))

	http.ListenAndServe(":3001", nil)
}

func LoginFunc(w http.ResponseWriter, r *http.Request) {
	loginTmpl := template.Must(template.ParseFiles("./frontend/templates/login.html"))
	userId := r.FormValue("userId")
	pwd := r.FormValue("pwd")
	fmt.Println(" userId ", userId, pwd, backend.Validate(userId, pwd))
	loginTmpl.Execute(w, "")
}

func QueryFunc(w http.ResponseWriter, r *http.Request) {
	homeTmpl := template.Must(template.ParseFiles("./frontend/templates/home.html"))
	IndexPage(r)
	logger.Log.Printf("Server v%s pid=%d started with processes: %d", VERSION, os.Getpid(), runtime.GOMAXPROCS(runtime.NumCPU()))
	homeTmpl.Execute(w, "")
}

func QueryWfFunc(w http.ResponseWriter, r *http.Request) {
	outputTmpl := template.Must(template.ParseFiles("./frontend/templates/output.html"))
	query := r.FormValue("query")
	fmt.Println(" query ", query)

	IndexPage(r)
	logger.Log.Printf("  Query Executed  ", query, "]")

	wfDbs := backend.GetWfDB("")
	noOfWorkflows := len(wfDbs)

	var wg sync.WaitGroup
	wfRespond := make(chan backend.Webdata, noOfWorkflows)
	wg.Add(noOfWorkflows)

	for _, col := range wfDbs {
		go backend.GetDbDataOnChan(wfRespond, &wg, query, "Workflow Data", col)
	}

	wg.Wait()
	close(wfRespond)

	outputTmpl.Execute(w, wfRespond)
}

func IndexPage(r *http.Request) {
	// get client IP address
	ip, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		logger.Log.Printf("userIP: [", r.RemoteAddr, "] is not IP:port")
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		logger.Log.Printf("userIP: [", r.RemoteAddr, " port ", port, "] is not IP:port")
		return
	}

	query := r.FormValue("query")
	logger.Log.Printf("userIP: [", r.RemoteAddr, " [ ", query, "]")

	// The user could acccess the web server via a proxy or load balancer.
	// The above IP address will be the IP address of the proxy or load balancer
	// and not the user's machine. Let's read the r.header "X-Forwarded-For (XFF)".
	// In our example the value returned is nil, we consider there is no proxy,
	// the IP indicates the user's address.
	// WARNING: this header is optional and will only be defined when site is
	// accessed via non-anonymous proxy and takes precedence over RemoteAddr.
	// (read https://tools.ietf.org/html/rfc7239 before any further use).
	proxied := r.Header.Get("X-FORWARDED-FOR")
	if proxied != "" {
		logger.Log.Printf("<p>Forwarded for: %s</p>", proxied)
	}
	logger.Log.Printf("<p>X-Forwarded-For : indicates this is the user's address (no proxy)<p>\n")
}
