package main

import (
	"io"
	"net/http"
	"log"
	"encoding/json"
	"os"
	"encoding/base64"
	"bytes"
	"golang.org/x/crypto/openpgp"
	"io/ioutil"
)


// https://golang.org/pkg/net/http/#ServeMux.Handle
// hello world, the web server

func DecryptSecret(s string) string {

        //all taken from here https://gist.github.com/stuart-warren/93750a142d3de4e8fdd2
	//TODO: Figure out how i should properly close this file.
        keyringFileBuffer, err := os.Open("./gpgkeys/secring.gpg")
	if err != nil {
		log.Println(err)
		return ""
	}
	entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		log.Println(err)
		return ""
	}
	//entity := entityList[0]
        dec, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Println(err)
		return ""
	}
        md, err := openpgp.ReadMessage(bytes.NewBuffer(dec), entityList, nil, nil)
	if err != nil {
		log.Println(err)
		return ""
	}
	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		log.Println(err)
		return ""
	}
	decStr := string(bytes)

        return decStr

}


func HandleRequest(w http.ResponseWriter, req *http.Request) {
	//Grab the contents of the request.  Make sure we have a valid object.  If not return error
	// Encrypt a message
	// =================
	// echo -n "this was encrypted" | gpg --batch --trust-model always --encrypt -r "Test Accounting <test@fakedomain.com>" --homedir ./gpgkeys/ | base64 -w 0

	// Decrypt that message
	// ====================
	// curl -H "Content-Type: application/json" -X POST -d '{"Key":"xyz","GpgContents":"hQEMA5siECpJxYMfAQf8D3+/slhHH47advUt+J8UZaQhP+fXiwR/R/GHalti3mFg6GknMIr7lKbxitGMsBFkTn+Qt1Mp0EuM+Z/suG1Jl+ppZ3WD5KjoBjZbDB8JYQw9aqPhQ/gl/M0GMwrZeyIUbHhtTUaEGZz53lcVWec6HItyGhEW3Lhv9MwciX4CvYp6aj3bJ6XLLfvTZe0qFBibMx8/VYX2gaJLZTC4kkTs7JXDVSzjGRKITnlVCUt+Qb8t8NkC/MzDH5xHZEts+KOy7TF5AQ7LBJ5mJL3X0WJtnnA7MSFFHABiT5/0Oa+SsH2aCv3AOb8l7FipVLxC3oJEbdEyp8eMMr54+pWDRRoreNJNAZMa2Fyt4wsDEbm6gsQzjjJ7e2Lj8osnDfiLalIEgczYHkkMpUaqNWlb+ErEUewKuMYahNg49op5RGu+8J+S+8cEJb2cWNdlLoi3kME="}"}' http://localhost:8080/
	json := json.NewDecoder(req.Body)

	type Request struct {
		Key string
		GpgContents string
	}

	//TODO: Validate that the json keys Key and GpgConents exist in the json else return 400.
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
		w.Header().Set("Content-Type", "text/plain")
		decrypted := DecryptSecret(r.GpgContents)
		if decrypted == "" {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Invalid Request")
		} else {
		io.WriteString(w, decrypted)
		}

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
