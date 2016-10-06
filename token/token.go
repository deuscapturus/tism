// token provides JWT validation and parsing.
package token

import (
	"../config"
	"../request"
	"../randid"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	"context"
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

func IsValid(w http.ResponseWriter, rc http.Request) (error, http.Request) {
	var req request.Request
	req = rc.Context().Value("request").(request.Request)

	token, err := parseToken(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		return err, rc
	}

	if token.Valid {
		return nil, rc
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		return errors.New("Token is not valid"), rc
	}
	return nil, rc
}

func GetAllowedKeys(w http.ResponseWriter, rc http.Request) (error, http.Request) {
	var req request.Request
	req = rc.Context().Value("request").(request.Request)

	token, err := parseToken(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		return err, rc
	}

	var claims []uint64
        for _, j := range token.Claims.(*JwtClaimsMap).Keys {
                j, err := strconv.ParseUint(j, 16, 64)
                if err != nil {
                        log.Println(err)
			return err, rc
                }
                claims = append(claims, j)
        }

	context := context.WithValue(rc.Context(), "claims", claims)
	return nil, *rc.WithContext(context)
}
// IssueJWT return a valid jwt with these statically defined scope values.
func IssueToken(w http.ResponseWriter, rc http.Request) (error, http.Request) {
	var req request.Request
	req = rc.Context().Value("request").(request.Request)

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
