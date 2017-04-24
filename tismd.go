// main package for the Immutable Secrets Managements API.  This service will decrypt and return a pgp encrypted messages if the requester has a valid jwt token that lists the pgp key's id in it's scope.
package main

import (
	"crypto/tls"
	"github.com/deuscapturus/tism/config"
	"github.com/deuscapturus/tism/encryption"
	"github.com/deuscapturus/tism/mytls"
	"github.com/deuscapturus/tism/randid"
	"github.com/deuscapturus/tism/request"
	"github.com/deuscapturus/tism/token"
	"log"
	"net/http"
)

func init() {
	// Generate admin token if requested
	if config.Config.GenAdminToken {
		adminToken, err := token.GenerateToken(
			[]string{"ALL"},
			99999999999,
			randid.Generate(32),
			1,
		)
		if err != nil {
			log.Println(err)
		}
		log.Println(adminToken)
	}
	// Generate new TLS cert if requested
	if config.Config.GenCert {
		mytls.Generate(config.Config.TLSDir)
	}
}

// main function.  Start http server and provide routes.
func main() {

	server := http.Server{
		Addr:         config.Config.Address + ":" + config.Config.Port,
		TLSConfig:    mytls.TLSConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	//Routes
	http.Handle("/decrypt", Handle(
		request.Parse,
		token.Parse,
		encryption.SetMyKeyRing,
		encryption.Decrypt,
	))
	http.Handle("/encrypt", Handle(
		request.Parse,
		token.Parse,
		encryption.SetMyKeyRing,
		encryption.Encrypt,
	))
	http.Handle("/token/info", Handle(
		request.Parse,
		token.Parse,
		token.Info,
	))
	http.Handle("/token/new", Handle(
		request.Parse,
		token.Parse,
		token.IsAdmin,
		token.New,
	))
	http.Handle("/key/new", Handle(
		request.Parse,
		token.Parse,
		token.IsAdmin,
		encryption.NewKey,
	))
	http.Handle("/key/list", Handle(
		request.Parse,
		token.Parse,
		encryption.SetMyKeyRing,
		encryption.ListKeys,
	))

	http.Handle("/key/get", Handle(
		request.Parse,
		token.Parse,
		encryption.SetMyKeyRing,
		encryption.GetKey,
	))

	http.Handle("/key/delete", Handle(
		request.Parse,
		token.Parse,
		token.IsAdmin,
		encryption.SetMyKeyRing,
		encryption.DeleteKey,
	))

	http.Handle("/", http.FileServer(http.Dir("./client")))

	log.Fatal(server.ListenAndServeTLS(
		config.Config.TLSDir+config.Config.TLSCertFile,
		config.Config.TLSDir+config.Config.TLSKeyFile,
	))
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
