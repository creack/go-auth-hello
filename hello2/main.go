package main

import (
	"fmt"
	"log"
	"net/http"
)

// hello is a dummy http endpoint replying "hello"
func hello(w http.ResponseWriter, req *http.Request) {
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
	fmt.Fprintf(w, "hello, %s\n", user)
}

func main() {
	http.HandleFunc("/", hello)
	fmt.Print("Ready on :8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error serving http: %s", err)
	}
}
