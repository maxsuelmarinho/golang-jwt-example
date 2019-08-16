package main

import (
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

var signingKey = []byte("myawesomesecret")

func main() {
	http.Handle("/", isAuthorized(homeHandler))
	log.Println("Running server on port 9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("home handler")
	fmt.Fprintf(w, "Hello World")	
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return signingKey, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}