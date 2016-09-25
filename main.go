// Package POC API will decrypt and return a pgp encrypted message if the requester has a valid jwt token that lists the pgp decrpytion id in it's scope.
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
	"golang.org/x/crypto/openpgp/armor"
	"io/ioutil"
	"github.com/dgrijalva/jwt-go"
	"time"
	"strconv"
)

// SigningKey to sign and validate the jwt
var SigningKey = []byte("sooFreakingsecret")

// DecryptSecret decrypt the given string.
// Do not return it if the decryption key is not in the k uint64 slice
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

// ValidateJWT validate string jwt and return true/false if valid along with a slice of uint64 pgp key id's.
func ValidateJWT(t string) (bool, []uint64) {
	type JwtClaimsMap struct {
	    Keys []string `json:"scopes"`
	    jwt.StandardClaims
	}
	//Validate jwt and return true or false plus a list of gpg private keys the requester has permission to.
	token, err := jwt.ParseWithClaims(t, &JwtClaimsMap{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SigningKey), nil
	})
	if err != nil {
		log.Println(err)
		return false, []uint64{0x000000}
	}
	if token.Valid {
		var claims []uint64
		for _, j := range token.Claims.(*JwtClaimsMap).Keys {
			j, err :=  strconv.ParseUint(j, 16, 64)
			if err != nil {
				log.Println(err)
			}
			claims = append(claims, j)
		}
		return true, claims
	} else {
		return false, []uint64{0x000000}
	}
}

func GetKeyRing() (entityList []*openpgp.Entity) {
	// Get PGP Encryption Keys
	keyringFileBuffer, err := os.Open("./gpgkeys/secring.gpg")
	if err != nil {
		log.Println(err)
		return
	}
	entityList, err = openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		log.Println(err)
		return
	}
	return

}
//Create ASscii Armor from openpgp.Entity
func PubEntToAsciiArmor(pubEnt *openpgp.Entity) (asciiEntity string) {
    gotWriter := bytes.NewBuffer(nil)
    wr, err := armor.Encode(gotWriter, openpgp.PublicKeyType, nil)
    if err != nil {
	log.Println(err)
	return
    }
    if pubEnt.Serialize(wr) != nil {
	log.Println(err)
    }
    if wr.Close() != nil {
	log.Println(err)
    }
    asciiEntity = gotWriter.String()
    return
}

// ListPGPKeys return json list of keys with metadata including id
func ListPGPKeys(w http.ResponseWriter, req *http.Request) {

	//TODO: Authenticate user maybe?

	var list []map[string]string
	JsonEncode := json.NewEncoder(w)
        // Return the id and pub key in json

	for _, entity := range GetKeyRing() {
		m := make(map[string]string)
		m["CreationTime"] = entity.PrimaryKey.CreationTime.String()
		m["Id"] = strconv.FormatUint(entity.PrimaryKey.KeyId, 16)
		for Name, _ := range entity.Identities {
			m["Name"] = Name
		}
		list = append(list, m)
	}

        w.Header().Set("Content-Type", "text/json")
        JsonEncode.Encode(list)
}

// GeneratePGPKey will create a new private/public gpg key pair
// and return the private key id and public key.
func GeneratePGPKey(w http.ResponseWriter, req *http.Request) {
/*
Create New PGP Key
==================
curl -H "Content-Type: application/json" -X POST -d '{"Key":"asdf", "Name":"testkey", "Comment":"test comment","Email":"test@key.test"}' http://localhost:8080/createpgp
*/
	JsonDecode := json.NewDecoder(req.Body)

	type RequestNewPGP struct {
		Key string
		Name string
		Comment string
		Email string
	}

	var r RequestNewPGP
	err := JsonDecode.Decode(&r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Invalid JSON in Request Body")
		log.Println(err)
		return
	}

	//TODO: Some authentication needed
	NewEntity, err := openpgp.NewEntity(r.Name, r.Comment, r.Email, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Unable to create PGP key based on input")
		log.Println(err)
		return
	}

	f, err := os.OpenFile("./gpgkeys/secring.gpg", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Println(err)
		return
	}

	NewEntity.SerializePrivate(f, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if f.Close() != nil {
		log.Println(err)
		return
	}

	type ReturnNewPGP struct {
		Id string
		PubKey string
	}

	// Return the id and pub key in json
	JsonEncode := json.NewEncoder(w)

	NewEntityId := strconv.FormatUint(NewEntity.PrimaryKey.KeyId, 16) 
	NewEntityPublicKey := PubEntToAsciiArmor(NewEntity)

	w.Header().Set("Content-Type", "text/json")
	p := &ReturnNewPGP{NewEntityId, NewEntityPublicKey}
	JsonEncode.Encode(*p)
}

// IssueJWT return a valid jwt with these statically defined scope values.
// This is only for testing.
func IssueJWT(w http.ResponseWriter, req *http.Request) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"scopes": []string{
			"4ceaf7b0e5e4040c",
			"4ceaf7b0e5e4040d",
			"8b3a8899b08ee628",
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

// HandleRequest the route for all decryption requests
func HandleRequest(w http.ResponseWriter, req *http.Request) {
/*
Encrypt a message
=================
echo -n "this was encrypted" | gpg --batch --trust-model always --encrypt -r "Test Accounting <test@fakedomain.com>" --homedir ./gpgkeys/ | base64 -w 0

Decrypt that message
====================
curl -H "Content-Type: application/json" -X POST -d '{"Key":"`curl -s http://localhost:8080/key`","GpgContents":"hQEMA5siECpJxYMfAQf8D3+/slhHH47advUt+J8UZaQhP+fXiwR/R/GHalti3mFg6GknMIr7lKbxitGMsBFkTn+Qt1Mp0EuM+Z/suG1Jl+ppZ3WD5KjoBjZbDB8JYQw9aqPhQ/gl/M0GMwrZeyIUbHhtTUaEGZz53lcVWec6HItyGhEW3Lhv9MwciX4CvYp6aj3bJ6XLLfvTZe0qFBibMx8/VYX2gaJLZTC4kkTs7JXDVSzjGRKITnlVCUt+Qb8t8NkC/MzDH5xHZEts+KOy7TF5AQ7LBJ5mJL3X0WJtnnA7MSFFHABiT5/0Oa+SsH2aCv3AOb8l7FipVLxC3oJEbdEyp8eMMr54+pWDRRoreNJNAZMa2Fyt4wsDEbm6gsQzjjJ7e2Lj8osnDfiLalIEgczYHkkMpUaqNWlb+ErEUewKuMYahNg49op5RGu+8J+S+8cEJb2cWNdlLoi3kME="}"}' http://localhost:8080/
*/
	json := json.NewDecoder(req.Body)

	type RequestDecrypt struct {
		Key string
		GpgContents string
	}

	var r RequestDecrypt
	err := json.Decode(&r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Invalid JSON in Request Body")
		log.Println(err)
		return
	}


	//Validate JWT
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
		// Return 404 if JWT is not valid
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Unauthorized")
		return
	}
}

// main function.  Start http server and provide routes.
func main() {
	server := http.Server{
		Addr: ":8080",
	}

	//Routes
	http.HandleFunc("/", HandleRequest)
	http.HandleFunc("/key", IssueJWT)
	http.HandleFunc("/createpgp", GeneratePGPKey)
	http.HandleFunc("/listpgp", ListPGPKeys)
	log.Fatal(server.ListenAndServe())
}
