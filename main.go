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

	"/backend"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const (
	VERSION     = "0.01"
	STATIC_PATH = "static"
)

func main() {
	// public views
	http.HandleFunc("/", LoginFunc)
	// private views
	http.HandleFunc("/query", QueryFunc)
	http.HandleFunc("/wf", QueryWfFunc)

	//http.HandleFunc("/query", PostOnly(BasicAuth(QueryFunc)))
	//http.HandleFunc("/wf", PostOnly(BasicAuth(QueryWfFunc)))

	//http.HandleFunc("/login", Login)

	fs := http.FileServer(http.Dir("static/"))
	http.Handle(STATIC_PATH, http.StripPrefix(STATIC_PATH, fs))

	http.ListenAndServe(":3000", nil)
}

func LoginFunc(w http.ResponseWriter, r *http.Request) {
	loginTmpl := template.Must(template.ParseFiles(STATIC_PATH + "/login.html"))
	userId := r.FormValue("uname")
	pwd := r.FormValue("pwd")

	backend.Log.Printf("UI userId " + userId + " pwd " + pwd)

	fmt.Println(" userId ", userId, pwd, backend.Validate(userId, pwd))
	loginTmpl.Execute(w, "")
}

func QueryFunc(w http.ResponseWriter, r *http.Request) {

	homeTmpl := template.Must(template.ParseFiles(STATIC_PATH + "/home.html"))
	IndexPage(r)
	backend.Log.Printf("Server v%s pid=%d started with processes: %d", VERSION, os.Getpid(), runtime.GOMAXPROCS(runtime.NumCPU()))
	homeTmpl.Execute(w, "")
}

func QueryWfFunc(w http.ResponseWriter, r *http.Request) {
	outputTmpl := template.Must(template.ParseFiles(STATIC_PATH + "/output.html"))
	query := r.FormValue("query")
	fmt.Println(" query ", query)

	IndexPage(r)
	backend.Log.Printf("  Query Executed  ", query, "]")

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
		backend.Log.Printf("userIP: [", r.RemoteAddr, "] is not IP:port")
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		backend.Log.Printf("userIP: [", r.RemoteAddr, " port ", port, "] is not IP:port")
		return
	}

	query := r.FormValue("query")
	backend.Log.Printf("userIP: [", r.RemoteAddr, " [ ", query, "]")

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
		backend.Log.Printf("<p>Forwarded for: %s</p>", proxied)
	}
	backend.Log.Printf("<p>X-Forwarded-For : indicates this is the user's address (no proxy)<p>\n")
}
