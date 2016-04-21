package main

import (
	"fmt"
	"log"
	"net/http"
)

// hello is a dummy http endpoint replying "hello"
func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "hello\n")
}

func main() {
	http.HandleFunc("/", hello)
	fmt.Print("Ready on :8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error serving http: %s", err)
	}
}
