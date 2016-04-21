package main

import (
	"fmt"
	"log"
	"net/http"
)

// Controller holds the state of the application.
type Controller struct {
	motd           string
	masterPassword string
}

// authHandler is our prototype for authenticated endpoints.
type authHandler func(http.ResponseWriter, *http.Request, string)

// hello is a dummy http endpoint replying the message of the day.
func (c *Controller) hello(w http.ResponseWriter, req *http.Request, user string) {
	fmt.Fprintf(w, "%s, %s\n", c.motd, user)
}

// MWAuth is a middleware checking for auth.
func (c *Controller) MWAuth(hdlr authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user, password, ok := req.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid auth settings\n")
			return
		}
		if password != c.masterPassword {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Invalid password\n")
			return
		}
		hdlr(w, req, user)
	}
}

func main() {
	c := &Controller{
		motd:           "hello",
		masterPassword: "world",
	}
	http.HandleFunc("/", c.MWAuth(c.hello))
	fmt.Print("Ready on :8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error serving http: %s", err)
	}
}
