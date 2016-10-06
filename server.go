// main package for the Immutable Secrets Managements API.  It will decrypt and return a pgp encrypted messages if the requester has a valid jwt token that lists the pgp decrpytion id in it's scope.
package main

import (
	"./token"
	"./config"
	"./request"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"errors"
)

var ConfigFilePath = string("./config.yaml")

type MyEntityList struct {
	openpgp.EntityList
}
var PGPKeyring = MyEntityList{}

// main function.  Start http server and provide routes.
func main() {

	config.Load(ConfigFilePath)
	PGPKeyring.GetKeyRing()

	server := http.Server{
		Addr: ":8080",
	}

	//Routes
	http.Handle("/decrypt", Handle(
		request.ParseRequest,
		token.IsValid,
		token.GetAllowedKeys,
		DecryptSecret,
	))
	http.Handle("/token/new", Handle(
		request.ParseRequest,
		token.IsValid,
		token.IssueToken,
	))
	http.Handle("/key/new", Handle(
		request.ParseRequest,
		token.IsValid,
		GeneratePGPKey,
	))
	http.Handle("/key/list", Handle(
		request.ParseRequest,
		token.IsValid,
		ListPGPKeys,
	))

	log.Fatal(server.ListenAndServe())
}

type Handler func (w http.ResponseWriter, rc http.Request) (error, http.Request)

// Handle is middleware chain handler
func Handle(handlers ...Handler) (http.Handler) {
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

// DecryptSecret decrypt the given string.
// Do not return it if the decryption key is not in the k uint64 slice
func DecryptSecret(w http.ResponseWriter, rc http.Request) (error, http.Request) {

	//var req request.Request
	req := rc.Context().Value("request").(request.Request)

	dec, err := base64.StdEncoding.DecodeString(req.GpgContents)
	if err != nil {
		log.Println(err)
		return err, rc
	}
	k := rc.Context().Value("claims").([]uint64)
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
			w.Write(message)
			return nil, rc
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
	return errors.New("Nothing Decrypted"), rc

}

// StringInSlice is self explanatory.  Return true or false.
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

	_, err := os.Stat(config.Config.KeyRingFilePath)
	var KeyringFileBuffer *os.File

	if os.IsNotExist(err) {
		KeyringFileBuffer, err = os.Create(config.Config.KeyRingFilePath)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		KeyringFileBuffer, err = os.Open(config.Config.KeyRingFilePath)
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
func ListPGPKeys(w http.ResponseWriter, rc http.Request) (error, http.Request) {

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
	return nil, rc
}

// GeneratePGPKey will create a new private/public gpg key pair
// and return the private key id and public key.
func GeneratePGPKey(w http.ResponseWriter, rc http.Request) (error, http.Request) {
        var req request.Request
	req = rc.Context().Value("request").(request.Request)

	//TODO: Some authentication needed
	NewEntity, err := openpgp.NewEntity(req.Name, req.Comment, req.Email, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		log.Println(err)
		return err, rc
	}

	f, err := os.OpenFile(config.Config.KeyRingFilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Println(err)
		return err, rc
	}
	defer f.Close()

	NewEntity.SerializePrivate(f, nil)
	if err != nil {
		log.Println(err)
		return err, rc
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
	return nil, rc
}
