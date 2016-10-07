// main package for the Immutable Secrets Managements API.  This service will decrypt and return a pgp encrypted messages if the requester has a valid jwt token that lists the pgp key's id in it's scope.
package main

import (
	"./config"
	"./encryption"
	"./request"
	"./token"
	"log"
	"net/http"
)

var ConfigFilePath = string("./config.yaml")

// main function.  Start http server and provide routes.
func main() {

	config.Load(ConfigFilePath)

	encryption.KeyRing.GetKeyRing()

	server := http.Server{
		Addr: ":8080",
	}

	//Routes
	http.Handle("/decrypt", Handle(
		request.ParseRequest,
		token.IsValid,
		token.GetAllowedKeys,
		encryption.Decrypt,
	))
	http.Handle("/token/new", Handle(
		request.ParseRequest,
		token.IsValid,
		token.IssueToken,
	))
	http.Handle("/key/new", Handle(
		request.ParseRequest,
		token.IsValid,
		encryption.NewKey,
	))
	http.Handle("/key/list", Handle(
		request.ParseRequest,
		token.IsValid,
		encryption.ListKeys,
	))
	http.Handle("/key/get", Handle(
		request.ParseRequest,
		token.IsValid,
		encryption.GetKey,
	))

	log.Fatal(server.ListenAndServe())
}

type Handler func(w http.ResponseWriter, rc http.Request) (error, http.Request)

// Handle is middleware chain handler
func Handle(handlers ...Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc := *r.WithContext(r.Context())
		var err error
		for _, handler := range handlers {
			err, rc = handler(w, rc)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
		}
	})
}
