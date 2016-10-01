/* Package POC API will decrypt and return a pgp encrypted message if the requester has a valid jwt token that lists the pgp decrpytion id in it's scope.

Create New PGP Key
==================
curl -H "Content-Type: application/json" -X POST -d '{"Key":"asdf", "Name":"testkey", "Comment":"test comment","Email":"test@key.test"}' http://localhost:8080/createpgp

Encrypt a Message
=================
echo -n "this was encrypted" | gpg --batch --trust-model always --encrypt -r "Test Accounting <test@fakedomain.com>" --homedir ./gpgkeys/ | base64 -w 0

List availble PGP Decryption Keys
=================================
curl http://localhost:8080/listpgp

Get a JWT Token
===============
curl -H "Content-Type: application/json" -X POST -d '{"Key":"123","Scope":["SOMEGPGID1","SOMEGPGID2"]' http://localhost:8080/key

Decrypt a Secret
================
curl -H "Content-Type: application/json" -X POST -d '{"Key":"`curl -s http://localhost:8080/key`","GpgContents":"hQEMA5siECpJxYMfAQf8D3+/slhHH47advUt+J8UZaQhP+fXiwR/R/GHalti3mFg6GknMIr7lKbxitGMsBFkTn+Qt1Mp0EuM+Z/suG1Jl+ppZ3WD5KjoBjZbDB8JYQw9aqPhQ/gl/M0GMwrZeyIUbHhtTUaEGZz53lcVWec6HItyGhEW3Lhv9MwciX4CvYp6aj3bJ6XLLfvTZe0qFBibMx8/VYX2gaJLZTC4kkTs7JXDVSzjGRKITnlVCUt+Qb8t8NkC/MzDH5xHZEts+KOy7TF5AQ7LBJ5mJL3X0WJtnnA7MSFFHABiT5/0Oa+SsH2aCv3AOb8l7FipVLxC3oJEbdEyp8eMMr54+pWDRRoreNJNAZMa2Fyt4wsDEbm6gsQzjjJ7e2Lj8osnDfiLalIEgczYHkkMpUaqNWlb+ErEUewKuMYahNg49op5RGu+8J+S+8cEJb2cWNdlLoi3kME="}"}' http://localhost:8080/

*/
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-yaml/yaml"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// SigningKey to sign and validate the jwt
var SigningKey = []byte("sooFreakingsecret")
var KeyRingFilePath = string("./gpgkeys/secring.gpg")
var ConfigFilePath = string("./config.yaml")

type Configuration struct {
	Port        int      `yaml:"port,omitempty"`
	RevokedJWTs []string `yaml:"revoked_api_keys,omitempty"`
}

type MyEntityList struct {
	openpgp.EntityList
}

var PGPKeyring = MyEntityList{}

var Config = Configuration{}

// main function.  Start http server and provide routes.
func main() {

	Config.LoadConfiguration(ConfigFilePath)
	PGPKeyring.GetKeyRing()

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

func (Config *Configuration) LoadConfiguration(ConfigPath string) {

	ConfigFileBytes, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		log.Println(err)
	}
	err = yaml.Unmarshal(ConfigFileBytes, &Config)
	if err != nil {
		log.Println(err)
	}
	return
}

// DecryptSecret decrypt the given string.
// Do not return it if the decryption key is not in the k uint64 slice
func DecryptSecret(s string, k []uint64) string {

	dec, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Println(err)
		return ""
	}

	//only load keys found in var k
	for _, keyid := range k {
		tryKeys := PGPKeyring.KeysById(keyid)
		for _, tryKey := range tryKeys {
			var TryKeyEntityList openpgp.EntityList

			TryKeyEntityList = append(TryKeyEntityList, tryKey.Entity)
			md, err := openpgp.ReadMessage(bytes.NewBuffer(dec), TryKeyEntityList, nil, nil)
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
		Keys  []string `json:"scopes"`
		JWTid string   `json:"jti"`
		jwt.StandardClaims
	}
	//Validate jwt and return true or false plus a list of gpg private keys the requester has permission to.
	token, err := jwt.ParseWithClaims(t, &JwtClaimsMap{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SigningKey), nil
	})
	if err != nil {
		log.Println(err)
		return false, []uint64{0}
	}

	if StringInSlice(token.Claims.(*JwtClaimsMap).JWTid, Config.RevokedJWTs) {
		log.Println("JWT Revoked")
		return false, []uint64{0}
	}

	if token.Valid {
		var claims []uint64
		for _, j := range token.Claims.(*JwtClaimsMap).Keys {
			j, err := strconv.ParseUint(j, 16, 64)
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

func StringInSlice(s string, slice []string) bool {
	for _, item := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetKeyRing return pgp keyring from a file location
func (PGPKeyring *MyEntityList) GetKeyRing() {

	_, err := os.Stat(KeyRingFilePath)
	var KeyringFileBuffer *os.File

	if os.IsNotExist(err) {
		KeyringFileBuffer, err = os.Create(KeyRingFilePath)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		KeyringFileBuffer, err = os.Open(KeyRingFilePath)
		if err != nil {
			log.Println(err)
			return
		}
	}

	EntityList, err := openpgp.ReadKeyRing(KeyringFileBuffer)
	if err != nil {
		log.Println(err)
		return
	}
	*PGPKeyring = MyEntityList{EntityList}

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

// ListPGPKeys return json list of keys with metadata including id.
func ListPGPKeys(w http.ResponseWriter, req *http.Request) {

	//TODO: Authenticate user maybe?
	var list []map[string]string
	JsonEncode := json.NewEncoder(w)
	// Return the id and pub key in json

	for _, entity := range PGPKeyring.EntityList {
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

	JsonDecode := json.NewDecoder(req.Body)

	type RequestNewPGP struct {
		Key     string
		Name    string
		Comment string
		Email   string
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

	f, err := os.OpenFile(KeyRingFilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	NewEntity.SerializePrivate(f, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Reload the Keyring after the new key is saved.
	defer PGPKeyring.GetKeyRing()

	type ReturnNewPGP struct {
		Id     string
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
func IssueJWT(w http.ResponseWriter, req *http.Request) {

	JsonDecode := json.NewDecoder(req.Body)

	type RequestNewJWT struct {
		Key   string
		Scope []string
	}

	var r RequestNewJWT
	err := JsonDecode.Decode(&r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Invalid JSON in Request Body")
		log.Println(err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"scopes": r.Scope,
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
		"jti":    "1234123412341234123412341234",
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

	json := json.NewDecoder(req.Body)

	type RequestDecrypt struct {
		Key         string
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
