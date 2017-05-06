// request package will parse all requests
package request

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// Request put fields for all requests in this struct.
type Request struct {
	Token     string   `json:"token"`
	Keys      []string `json:"keys"`
	Admin     int      `json:"admin"`
	EncSecret string   `json:"encsecret"`
	Encoding  string   `json:"encoding"`
	DecSecret string   `json:"decsecret"`
	Key       string   `json:"key"`
	Name      string   `json:"name"`
	Comment   string   `json:"comment"`
	Email     string   `json:"email"`
	Id        string   `json:"id"`
}

// Parse get the json message body and map it to our Request struct.
// This is the only way to get information from the request body through the middlware chain.  The message body is destroyed after the first read here.
func Parse(w http.ResponseWriter, rc http.Request) (error, http.Request) {

	var req = Request{}

	json := json.NewDecoder(rc.Body)

	err := json.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		return errors.New("Invalid JSON request in body"), rc
	}
	context := context.WithValue(rc.Context(), "request", req)
	return nil, *rc.WithContext(context)
}
