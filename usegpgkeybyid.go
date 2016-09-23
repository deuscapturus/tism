package main

import (
	"log"
	"os"
	"encoding/base64"
	"bytes"
	"golang.org/x/crypto/openpgp"
	"io/ioutil"
	"reflect"
)


var SigningKey = []byte("sooFreakingsecret")
var GpgContents = string("hQEMA5siECpJxYMfAQf8D3+/slhHH47advUt+J8UZaQhP+fXiwR/R/GHalti3mFg6GknMIr7lKbxitGMsBFkTn+Qt1Mp0EuM+Z/suG1Jl+ppZ3WD5KjoBjZbDB8JYQw9aqPhQ/gl/M0GMwrZeyIUbHhtTUaEGZz53lcVWec6HItyGhEW3Lhv9MwciX4CvYp6aj3bJ6XLLfvTZe0qFBibMx8/VYX2gaJLZTC4kkTs7JXDVSzjGRKITnlVCUt+Qb8t8NkC/MzDH5xHZEts+KOy7TF5AQ7LBJ5mJL3X0WJtnnA7MSFFHABiT5/0Oa+SsH2aCv3AOb8l7FipVLxC3oJEbdEyp8eMMr54+pWDRRoreNJNAZMa2Fyt4wsDEbm6gsQzjjJ7e2Lj8osnDfiLalIEgczYHkkMpUaqNWlb+ErEUewKuMYahNg49op5RGu+8J+S+8cEJb2cWNdlLoi3kME=")

func DecryptSecret(s string) string {

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
	
	md, err := openpgp.ReadMessage(bytes.NewBuffer(dec), entityList, nil, nil)
	if err != nil {
		log.Println(err)
		return ""
	}
	messagebytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		log.Println(err)
		return ""
	}
	log.Println(entityList[0].PrimaryKey.KeyId)
	log.Println(reflect.TypeOf(entityList[0].PrimaryKey.KeyId))


	var tryKeyentityList openpgp.EntityList
	tryKeyentityList = append(tryKeyentityList, entityList.KeysById(entityList[0].PrimaryKey.KeyId)[0].Entity)
	newmd, err := openpgp.ReadMessage(bytes.NewBuffer(dec), tryKeyentityList, nil, nil)
	if err != nil {
		log.Println(err)
		return ""
	}
	newmessagebytes, err := ioutil.ReadAll(newmd.UnverifiedBody)
	if err != nil {
		log.Println(err)
		return ""
	}
	log.Println(string(newmessagebytes))

	return string(messagebytes)

}


func main() {
	log.Println(DecryptSecret(GpgContents))
}

