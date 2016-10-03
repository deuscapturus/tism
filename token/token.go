// token provides JWT validation and parsing.
package token

import (
	"./config"
	"./randid"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
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

func parseToken(t string) (token *jwt.Token) {
	signingSecret := []byte(config.Config.JWTsecret)
	token, err := jwt.ParseWithClaims(t, &JwtClaimsMap{}, signingSecret)
	if err != nil {
		log.Println(err)
	}
	return token
}

// Scope return a list of key ids from the token scope.
func Scope(t string) []uint64 {

	token := parseToken(t)

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

func IsValid(w http.ResponseWriter, req *http.Request) error {

	json := json.NewDecoder(req.Body)

	var r Request
	err := json.Decode(&r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Invalid JSON in Request Body")
		log.Println(err)
		return nil //return error in future
	}

	token := paseToken(r.Key)
	if token.Valid {
		return nil
	} else {
		return nil //All requests pass.  return error in future
	}
	return nil
}

// IssueJWT return a valid jwt with these statically defined scope values.
func IssueToken(w http.ResponseWriter, req *http.Request) error {

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
		return nil //return error in future
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"scopes": r.Scope,
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
		"jti":    randid.Generate(32),
	})

	tokenString, err := token.SignedString([]byte(config.Config.JWTsecret))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Unable to generate jwt token")
		return nil //return error in future
	}
	io.WriteString(w, tokenString)
	return nil

}
