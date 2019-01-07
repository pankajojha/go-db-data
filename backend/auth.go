package backend

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type handler func(w http.ResponseWriter, r *http.Request)

func BasicAuth(pass handler) handler {

	return func(w http.ResponseWriter, r *http.Request) {

		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !Validate(pair[0], pair[1]) {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		pass(w, r)
	}
}

//Validate ...
func Validate(username, password string) bool {

	Log.Printf("UserName ", username, " password ", password)

	if username == Config.UserName && password == "test" {

		return true
	}
	return false
}
