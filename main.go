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
	"github.com/dgrijalva/jwt-go"
	"time"
)


var SigningKey = []byte("sooFreakingsecret")
// https://golang.org/pkg/net/http/#ServeMux.Handle
// hello world, the web server

func IsKeyIDinIDSlice(id uint64, idslice []uint64) bool {
	for _, i := range idslice {
		if i == id {
			return true
		}
	}
	return false
}

func DecryptSecret(s string, k []uint64) string {

	//all taken from here https://gist.github.com/stuart-warren/93750a142d3de4e8fdd2
	//TODO: Figure out how i should properly close this file.
	keyringFileBuffer, err := os.Open("./gpgkeys/secring.gpg")
	if err != nil {
		log.Println(err)
		return ""
	}
        dec, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Println(err)
		return ""
	}
	entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		log.Println(err)
		return ""
	}
	//only load keys found in var k
	for _, keyid := range k {
		tryKeys := entityList.KeysById(keyid)
		for _, tryKey := range tryKeys {
			var tryKeyentityList openpgp.EntityList
			tryKeyentityList = append(tryKeyentityList, tryKey.Entity)
			md, err := openpgp.ReadMessage(bytes.NewBuffer(dec), tryKeyentityList, nil, nil)
			if err != nil {
				log.Println(err)
				continue
			}
			message, err := ioutil.ReadAll(md.UnverifiedBody)
			if err != nil {
				log.Println(err)
				continue
			}
			return string(message)
		}
	}
	return ""


}

func ValidateJWT(t string) (bool, []uint64) {
	//Validate jwt and return true or false plus a list of gpg private keys the requester has permission to.
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte(SigningKey), nil
	})
	if err != nil {
		log.Println(err)
		return false, []uint64{0x000000}
	}
	if token.Valid {
		return true, []uint64{0x4ceaf7b0e5e4040c,0x4ceaf7b0e5e4040d,}
	} else {
		return false, []uint64{0x000000}
	}
}

func IssueJWT(w http.ResponseWriter, req *http.Request) {
	//{
	//  "scopes": [
	//      "IDENTIFIEROFKEY1",
	//      "IDENTIFIEROFKEY2",
	//  ],
	//  "exp": 1446056652,
	//  "jti": "7e9c6a991f5a227fb7ebaa522536ae4c"
	//}
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	// curl http://localhost:8080/key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"scopes": []uint64{
			0x4ceaf7b0e5e4040c,
			0x4ceaf7b0e5e4040d,
		},
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"jti": "1234123412341234123412341234",
	})

	tokenString, err := token.SignedString(SigningKey)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Unable to generate jwt token")
	}
	io.WriteString(w, tokenString)

}

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	//Grab the contents of the request.  Make sure we have a valid object.  If not return error
	// Encrypt a message
	// =================
	// echo -n "this was encrypted" | gpg --batch --trust-model always --encrypt -r "Test Accounting <test@fakedomain.com>" --homedir ./gpgkeys/ | base64 -w 0

	// Decrypt that message
	// ====================
	// curl -H "Content-Type: application/json" -X POST -d '{"Key":"`curl http://localhost:8080/key`","GpgContents":"hQEMA5siECpJxYMfAQf8D3+/slhHH47advUt+J8UZaQhP+fXiwR/R/GHalti3mFg6GknMIr7lKbxitGMsBFkTn+Qt1Mp0EuM+Z/suG1Jl+ppZ3WD5KjoBjZbDB8JYQw9aqPhQ/gl/M0GMwrZeyIUbHhtTUaEGZz53lcVWec6HItyGhEW3Lhv9MwciX4CvYp6aj3bJ6XLLfvTZe0qFBibMx8/VYX2gaJLZTC4kkTs7JXDVSzjGRKITnlVCUt+Qb8t8NkC/MzDH5xHZEts+KOy7TF5AQ7LBJ5mJL3X0WJtnnA7MSFFHABiT5/0Oa+SsH2aCv3AOb8l7FipVLxC3oJEbdEyp8eMMr54+pWDRRoreNJNAZMa2Fyt4wsDEbm6gsQzjjJ7e2Lj8osnDfiLalIEgczYHkkMpUaqNWlb+ErEUewKuMYahNg49op5RGu+8J+S+8cEJb2cWNdlLoi3kME="}"}' http://localhost:8080/
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
		io.WriteString(w, "Invalid JSON in Request Body")
		log.Println(err)
		return
	}


	//if r.Key is a valid jwt decrypt.
	//Validate Key
	jwtvalid, gpgkeys := ValidateJWT(r.Key)
	if jwtvalid {
		w.Header().Set("Content-Type", "text/plain")
		decrypted := DecryptSecret(r.GpgContents, gpgkeys)
		if decrypted == "" {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Unable to Decrypt Message")
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
	http.HandleFunc("/key", IssueJWT)

	log.Fatal(server.ListenAndServe())
}
