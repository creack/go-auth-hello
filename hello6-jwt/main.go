package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// authHandler is our prototype for authenticated endpoints.
type authHandler func(http.ResponseWriter, *http.Request, string)

// Controller holds the state of the application.
type Controller struct {
	motd           string
	masterPassword string
	privateKey     []byte
	publicKey      []byte
}

// NewController loads the RSA keys and instantiate the Controller.
func NewController() (*Controller, error) {
	privateKey, err := ioutil.ReadFile("private.key")
	if err != nil {
		return nil, err
	}
	publicKey, err := ioutil.ReadFile("public.key")
	if err != nil {
		return nil, err
	}

	c := &Controller{
		motd:           "hello",
		masterPassword: "world",
		privateKey:     privateKey,
		publicKey:      publicKey,
	}
	return c, nil
}

// hello is a dummy http endpoint replying the message of the day.
func (c *Controller) hello(w http.ResponseWriter, req *http.Request, user string) {
	fmt.Fprintf(w, "%s, %s\n", c.motd, user)
}

// MWAuth is a middleware checking for auth.
func (c *Controller) MWAuth(hdlr authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Parse the token.
		token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
			return c.publicKey, nil
		})
		// If then token is plain invalid or expired, stop here.
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized\n")
			return
		}

		hdlr(w, req, token.Claims["user"].(string))
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

	signer := jwt.New(jwt.GetSigningMethod("RS256"))
	signer.Claims["user"] = user
	signer.Claims["exp"] = time.Now().UTC().Add(30 * time.Second).Unix()

	token, err := signer.SignedString(c.privateKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error generating new token\n")
		log.Printf("Error generating new token: %s", err)
		return
	}

	fmt.Fprintf(w, "%s\n", token)
}

func main() {
	c, err := NewController()
	if err != nil {
		log.Fatalf("Error instantiating the controller: %s", err)
	}

	http.HandleFunc("/", c.MWAuth(c.hello))
	http.HandleFunc("/login", c.Login)
	fmt.Print("Ready on :8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error serving http: %s", err)
	}
}
