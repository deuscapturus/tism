package main

import (
	"io"
	"net/http"
	"log"
	"encoding/json"
)

// https://golang.org/pkg/net/http/#ServeMux.Handle
// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	//Grab the contents of the request.  Make sure we have a valid object.  If not return error
	contents := json.NewDecoder(req.Body)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "hello, world!\n")
}


func main() {
	server := http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/", HelloServer)

	log.Fatal(server.ListenAndServe())
}
