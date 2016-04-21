package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Controller holds the state of the application.
type Controller struct {
	motd           string
	masterPassword string
	tokens         map[string]string
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
		auth := req.Header.Get("Authorization")

		const prefix = "Bearer "
		if !strings.HasPrefix(auth, prefix) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid auth header\n")
			return
		}
		user, ok := c.tokens[auth[len(prefix):]]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Invalid token\n")
			return
		}

		hdlr(w, req, user)
	}
}

// Login is our login endpoint.
func (c *Controller) Login(w http.ResponseWriter, req *http.Request) {
	// Validate the login.
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

	// Generate, store and return token.
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error generating random string for token\n")
		return
	}
	tokenStr := hex.EncodeToString(token)

	c.tokens[tokenStr] = user
	fmt.Fprintf(w, "%s\n", tokenStr)
}

func main() {
	c := &Controller{
		motd:           "hello",
		masterPassword: "world",
		tokens:         map[string]string{},
	}
	http.HandleFunc("/", c.MWAuth(c.hello))
	http.HandleFunc("/login", c.Login)
	fmt.Print("Ready on :8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error serving http: %s", err)
	}
}
