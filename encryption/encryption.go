// encrpytion package holds all pgp related tasks
package encryption

import (
	"../config"
	"../request"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type MyEntityList struct {
	openpgp.EntityList
}

type PublicKey struct {
	Id     string
	PubKey string
}

var KeyRing = MyEntityList{}

// Decrypt decrypt the given string.
// Do not return it if the decryption key is not in the k uint64 slice
func Decrypt(w http.ResponseWriter, rc http.Request) (error, http.Request) {

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
		tryKeys := KeyRing.KeysById(keyid)
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

// ListKeys return json list of keys with metadata including id.
func ListKeys(w http.ResponseWriter, rc http.Request) (error, http.Request) {

	var list []map[string]string
	JsonEncode := json.NewEncoder(w)

	for _, entity := range KeyRing.EntityList {
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

func GetKey(w http.ResponseWriter, rc http.Request) (error, http.Request) {

	var req request.Request
	req = rc.Context().Value("request").(request.Request)
	EntityId, err := strconv.ParseUint(req.Id, 16, 64)
	if err != nil {
		return err, rc
	}

	ThisKey := KeyRing.KeysById(EntityId)

	JsonEncode := json.NewEncoder(w)

	w.Header().Set("Content-Type", "text/json")
	// TODO: KeysById returns a slice, though there should only ever be one id per key.  For now assume only one key is ever returned.  Re-consider in the future.
	p := &PublicKey{req.Id, PubEntToAsciiArmor(ThisKey[0].Entity)}
	JsonEncode.Encode(*p)
	return nil, rc

}

// NewKey will create a new private/public gpg key pair
// and return the private key id and public key.
func NewKey(w http.ResponseWriter, rc http.Request) (error, http.Request) {
	var req request.Request
	req = rc.Context().Value("request").(request.Request)

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
	defer KeyRing.GetKeyRing()

	// Return the id and pub key in json
	JsonEncode := json.NewEncoder(w)

	NewEntityId := strconv.FormatUint(NewEntity.PrimaryKey.KeyId, 16)
	NewEntityPublicKey := PubEntToAsciiArmor(NewEntity)

	w.Header().Set("Content-Type", "text/json")
	p := &PublicKey{NewEntityId, NewEntityPublicKey}
	JsonEncode.Encode(*p)
	return nil, rc
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
func (KeyRing *MyEntityList) GetKeyRing() {

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
	*KeyRing = MyEntityList{EntityList}

	return
}

//PubEntToAsciiArmor create Ascii Armor from openpgp.Entity
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
