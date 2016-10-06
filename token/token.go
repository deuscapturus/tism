// token provides JWT validation and parsing.
package token

import (
	"../config"
	"../randid"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	"io/ioutil"
	"bytes"
//	"context"
)

type JwtClaimsMap struct {
	Keys  []string `json:"scopes"`
	JWTid string   `json:"jti"`
	jwt.StandardClaims
}

// ValidateJWT validate string jwt and return true/false if valid along with a slice of uint64 pgp key id's.

type RequestDecrypt struct {
	Token string `json:"token"`
}

func parseToken(t string) (token *jwt.Token, err error) {
	signingSecret := func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JWTsecret), nil
	}
	token, err = jwt.ParseWithClaims(t, &JwtClaimsMap{}, signingSecret)
	return token, err
}

// Scope return a list of key ids from the token scope.
func Scope(t string) []uint64 {

	token, err := parseToken(t)
	if err != nil {
		log.Println(err)
	}

	var claims []uint64

	for _, j := range token.Claims.(*JwtClaimsMap).Keys {
		j, err := strconv.ParseUint(j, 16, 64)
		if err != nil {
			log.Println(err)
		}
		claims = append(claims, j)
	}
	return claims
}
const Mytest = "fake"
func IsValid(w http.ResponseWriter, rc http.Request) (error, http.Request) {
	log.Println(rc.Context().Value(Mytest))
	type Request struct {
		Token string `json:"token"`
		Scope []string `json:"Scope"`
	}
	json := json.NewDecoder(rc.Body)

	var req Request
	err := json.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Invalid JSON in Request Body in IsValid")
		return err, rc
	}

	token, err := parseToken(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		return err, rc
	}

	if token.Valid {
		log.Println("Token is Valid")
		return nil, rc
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		return errors.New("Token is not valid"), rc //All requests pass.  return error in future
	}
	return nil, rc
}

// IssueJWT return a valid jwt with these statically defined scope values.
func IssueToken(w http.ResponseWriter, rc http.Request) (error, http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc.Body)
	log.Println("logging out req buffer in IssueToken")
	log.Println(buf)
	log.Println("finished logging out req buffer in IssueToken")

	reqbody, _ := ioutil.ReadAll(rc.Body)
	log.Println(string(reqbody))
	json := json.NewDecoder(rc.Body)

	type RequestNewToken struct {
		Token string `json:"token"`
		Scope []string
	}

	var req RequestNewToken
	err := json.Decode(&req)
	if err != nil {
		log.Println("in JsonDecode.Decode")
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return errors.New("Invalid JSON for some dumbass reason"), rc
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"scopes": req.Scope,
		"exp":    time.Now().Add(time.Hour * 30303).Unix(),
		"jti":    randid.Generate(32),
	})

	tokenString, err := token.SignedString([]byte(config.Config.JWTsecret))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return err, rc
	}
	io.WriteString(w, tokenString)
	return nil, rc

}
