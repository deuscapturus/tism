package main

import (
	"io"
	"net/http"
	"log"
	"encoding/json"
)


// https://golang.org/pkg/net/http/#ServeMux.Handle
// hello world, the web server

func DecryptSecret(s string) string {
	return s
}

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	//Grab the contents of the request.  Make sure we have a valid object.  If not return error
	// curl -H "Content-Type: application/json" -X POST -d '{"Key":"xyz","GpgContents":"xyz"}' http://localhost:8080/
	json := json.NewDecoder(req.Body)

	type Request struct {
		Key string
		GpgContents string
	}

	//TODO: if json fails return error with code 400
	var r Request
	err := json.Decode(&r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Invalid Request")
		log.Println(err)
		return
	}


	//if r.Key is a valid jwt decrypt.
	if r.Key == "xyz" {
		io.WriteString(w, DecryptSecret(r.Key))

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Unauthorized")
		return
		//return error code 401
	}
}


func main() {
	server := http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/", HandleRequest)

	log.Fatal(server.ListenAndServe())
}
