package main

import (
	"fmt"
	"log"
	"net/http"
)

// authHandler is our prototype for authenticated endpoints.
type authHandler func(http.ResponseWriter, *http.Request, string)

// hello is a dummy http endpoint replying "hello"
func hello(w http.ResponseWriter, req *http.Request, user string) {
	fmt.Fprintf(w, "hello, %s\n", user)
}

// MWAuth is a middleware checking for auth.
func MWAuth(hdlr authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user, password, ok := req.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid auth settings\n")
			return
		}
		if password != "world" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Invalid password\n")
			return
		}
		hdlr(w, req, user)
	}
}

func main() {
	http.HandleFunc("/", MWAuth(hello))
	fmt.Print("Ready on :8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error serving http: %s", err)
	}
}
