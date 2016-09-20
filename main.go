package main

import (
	"io"
	"net/http"
	"log"
	"encoding/json"
)


// https://golang.org/pkg/net/http/#ServeMux.Handle
// hello world, the web server
func HandleRequest(w http.ResponseWriter, req *http.Request) {
	//Grab the contents of the request.  Make sure we have a valid object.  If not return error
	json := json.NewDecoder(req.Body)

	type Request struct {
		Key string
		GpgContents string
	}

	var r Request
	decodedjson := json.Decode(&r)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, decodedjson.Key)
}


func main() {
	server := http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/", HandleRequest)

	log.Fatal(server.ListenAndServe())
}
