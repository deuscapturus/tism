// token provides JWT validation and parsing.
package token

import (
	"../config"
	"../randid"
	"../request"
	"../utils"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io"
	"log"
	"net/http"
	"time"
)

type JwtClaimsMap struct {
	Keys  []string `json:"keys"`
	Admin int      `json:"admin"`
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
	if err != nil {
		log.Println(err)
	}
	return token, err
}

func Parse(w http.ResponseWriter, rc http.Request) (error, http.Request) {
	var req request.Request
	req = rc.Context().Value("request").(request.Request)
	token, err := parseToken(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		return err, rc
	}

	if token.Valid {
		// set scope to string "ALL" in request context if requester has privilege to all keys.
		// else, set scope to slice of uint64 key ids from the token.
		var context1 context.Context
		if token.Claims.(*JwtClaimsMap).Keys[0] == "ALL" {
			context1 = context.WithValue(rc.Context(), "claims", "ALL")
		} else {
			claims := token.Claims.(*JwtClaimsMap).Keys
			context1 = context.WithValue(rc.Context(), "claims", claims)
		}

		var admin int
		admin = token.Claims.(*JwtClaimsMap).Admin
		context2 := context.WithValue(context1, "admin", admin)

		return nil, *rc.WithContext(context2)
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "text/plain")
	return errors.New("Token is not valid"), rc
}

// IsAdmin Only continue if the requestor is an admin.
func IsAdmin(w http.ResponseWriter, rc http.Request) (error, http.Request) {

	admin := rc.Context().Value("admin").(int)
	if admin >= 1 {
		return nil, rc
	}
	return errors.New("Requestor is not admin"), rc
}

// IssueJWT return a valid jwt with these statically defined scope values.
func New(w http.ResponseWriter, rc http.Request) (error, http.Request) {

	var req request.Request
	req = rc.Context().Value("request").(request.Request)

	authKeys := rc.Context().Value("claims")
	if rc.Context().Value("admin").(int) >= 0 {
		switch authKeys.(type) {
		case string:
			if authKeys.(string) != "ALL" {
				log.Println("A valid token is not configured correctly.")
				return errors.New("Permission Denied.  Requested Keys are not in requestors allowed scope"), rc
			}
		case []string:
			if !utils.AllStringsInSlice(req.Keys, rc.Context().Value("claims").([]string)) {
				log.Println("Requested Keys are not in requestors allowed")
				return errors.New("Permission Denied.  Requested Keys are not in requestors allowed scope"), rc
			}
		}
	} else {
		return errors.New("Permission Denied.  Not an admin token"), rc
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"keys":  req.Keys,
		"exp":   time.Now().Add(time.Hour * 30303).Unix(),
		"jti":   randid.Generate(32),
		"admin": req.Admin,
	})

	tokenString, err := token.SignedString([]byte(config.Config.JWTsecret))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return err, rc
	}

	io.WriteString(w, tokenString)
	return nil, rc

	return errors.New("Permission Denied"), rc
}
