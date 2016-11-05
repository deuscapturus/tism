// encrpytion package holds all pgp related tasks
package encryption

import (
	"bytes"
	"context"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/deuscapturus/tism/config"
	"github.com/deuscapturus/tism/request"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
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
	Id     string `json:"id"`
	PubKey string `json:"pubkey"`
}

var KeyRing = MyEntityList{}

func init() {
	KeyRing.GetKeyRing()
}

func SetMyKeyRing(w http.ResponseWriter, rc http.Request) (error, http.Request) {

	var MyKeyRing openpgp.EntityList

	AuthorizedKeys := rc.Context().Value("claims")

	switch AuthorizedKeys.(type) {
	case string:
		if AuthorizedKeys.(string) == "ALL" {
			MyKeyRing = KeyRing.EntityList
		}
	case []string:
		// Assemble a new entity list based on the outcome of KeysById
		keys := AuthorizedKeys.([]string)
		keysUint64 := stringsToUint64(keys)
		for _, keyid := range keysUint64 {
			for _, thisk := range KeyRing.KeysById(keyid) {
				MyKeyRing = append(MyKeyRing, thisk.Entity)
			}
		}
	}

	context := context.WithValue(rc.Context(), "MyKeyRing", MyKeyRing)
	return nil, *rc.WithContext(context)
}

// Decrypt decrypt the given string.
func Decrypt(w http.ResponseWriter, rc http.Request) (error, http.Request) {

	req := rc.Context().Value("request").(request.Request)
	MyKeyRing := rc.Context().Value("MyKeyRing").(openpgp.EntityList)

	dec, err := base64.StdEncoding.DecodeString(req.EncSecret)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return err, rc
	}

	md, err := openpgp.ReadMessage(bytes.NewBuffer(dec), MyKeyRing, nil, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return err, rc
	}

	message, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return err, rc
	}

	w.Write(message)
	return nil, rc
}

// Encrypt encrypt the given string.
func Encrypt(w http.ResponseWriter, rc http.Request) (error, http.Request) {

	req := rc.Context().Value("request").(request.Request)
	MyKeyRing := rc.Context().Value("MyKeyRing").(openpgp.EntityList)

	EntityId := stringToUint64(req.Id)
	ThisKey := MyKeyRing.KeysById(EntityId)
	if len(ThisKey) == 0 {
		return errors.New("Key Id not found in requestors keyring"), rc
	}

	var ThisKeyEntity []*openpgp.Entity
	ThisKeyEntity = append(ThisKeyEntity, ThisKey[0].Entity)

	buf := new(bytes.Buffer)
	EncWriter, err := openpgp.Encrypt(buf, ThisKeyEntity, nil, nil, nil)
	if err != nil {
		return err, rc
	}

	_, err = EncWriter.Write([]byte(req.DecSecret))
	if err != nil {
		return err, rc
	}

	err = EncWriter.Close()
	if err != nil {
		return err, rc
	}

	encSecret, err := ioutil.ReadAll(buf)
	if err != nil {
		return err, rc
	}

	base64encSecret := base64.StdEncoding.EncodeToString(encSecret)

	w.Header().Set("Content-Type", "text/text")
	w.Write([]byte(base64encSecret))

	return nil, rc
}

// ListKeys return json list of keys with metadata including id.
func ListKeys(w http.ResponseWriter, rc http.Request) (error, http.Request) {

	var list []map[string]string
	JsonEncode := json.NewEncoder(w)

	MyKeyRing := rc.Context().Value("MyKeyRing").(openpgp.EntityList)

	if rc.Context().Value("claims").(string) == "ALL" {
		all := make(map[string]string)
		all["Id"] = "ALL"
		list = append(list, all)
	}

	for _, entity := range MyKeyRing {
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

	MyKeyRing := rc.Context().Value("MyKeyRing").(openpgp.EntityList)

	var req request.Request
	req = rc.Context().Value("request").(request.Request)
	EntityId, err := strconv.ParseUint(req.Id, 16, 64)
	if err != nil {
		return err, rc
	}

	ThisKey := MyKeyRing.KeysById(EntityId)
	if ThisKey == nil {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("Key not found"), rc
	}

	JsonEncode := json.NewEncoder(w)

	w.Header().Set("Content-Type", "text/json")
	// TODO: KeysById returns a slice, though there should only ever be one id per key.  For now assume only one key is ever returned.  Re-consider in the future.
	p := &PublicKey{req.Id, pubEntToAsciiArmor(ThisKey[0].Entity)}
	JsonEncode.Encode(*p)
	return nil, rc

}

// NewKey will create a new private/public gpg key pair
// and return the private key id and public key.
func NewKey(w http.ResponseWriter, rc http.Request) (error, http.Request) {
	var req request.Request
	req = rc.Context().Value("request").(request.Request)

	pgpConfig := &packet.Config{
		DefaultHash : crypto.SHA256,
	}

	NewEntity, err := openpgp.NewEntity(req.Name, req.Comment, req.Email, pgpConfig)
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
	NewEntityPublicKey := pubEntToAsciiArmor(NewEntity)

	w.Header().Set("Content-Type", "text/json")
	p := &PublicKey{NewEntityId, NewEntityPublicKey}
	JsonEncode.Encode(*p)
	return nil, rc
}

// stringsToUint64 convert a slice of strings to uint64.
func stringsToUint64(s []string) []uint64 {

	var uint64List []uint64

	for _, j := range s {
		j, err := strconv.ParseUint(j, 16, 64)
		if err != nil {
			log.Println(err)
		}
		uint64List = append(uint64List, j)
	}
	return uint64List
}

// stringsToUint64 convert string to uint64.
func stringToUint64(s string) uint64 {

	suint64, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		log.Println(err)
	}
	return suint64
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

// pubEntToAsciiArmor create Ascii Armor from openpgp.Entity
func pubEntToAsciiArmor(pubEnt *openpgp.Entity) (asciiEntity string) {

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
