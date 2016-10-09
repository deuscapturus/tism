// main package for the Immutable Secrets Managements API.  This service will decrypt and return a pgp encrypted messages if the requester has a valid jwt token that lists the pgp key's id in it's scope.
package main

import (
	"./config"
	"./randid"
	"./encryption"
	"./request"
	"./token"
	"log"
	"net/http"
)

// main function.  Start http server and provide routes.
func main() {

	config.Load()

	encryption.KeyRing.GetKeyRing()

	if config.Config.GenAdminToken {
		adminToken, err := token.GenerateToken([]string{"ALL"}, 99999999999, randid.Generate(32), 1)
		if err != nil {
			log.Println(err)
		}
		log.Println(adminToken)
	}


	server := http.Server{
		Addr: ":8080",
	}

	//Routes
	http.Handle("/decrypt", Handle(
		request.ParseRequest,
		token.Parse,
		encryption.SetMyKeyRing,
		encryption.Decrypt,
	))
	http.Handle("/token/new", Handle(
		request.ParseRequest,
		token.Parse,
		token.IsAdmin,
		token.New,
	))
	http.Handle("/key/new", Handle(
		request.ParseRequest,
		token.Parse,
		token.IsAdmin,
		encryption.NewKey,
	))
	http.Handle("/key/list", Handle(
		request.ParseRequest,
		token.Parse,
		encryption.SetMyKeyRing,
		encryption.ListKeys,
	))

	http.Handle("/key/get", Handle(
		request.ParseRequest,
		token.Parse,
		encryption.SetMyKeyRing,
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
