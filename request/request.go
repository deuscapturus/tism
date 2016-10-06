// request package will parse all requests
package request

import (
	"encoding/json"
	"log"
	"net/http"
	"errors"
	"context"
)

type Request struct {
	Token string `json:"token"`
	Scope []string `json:"Scope"`
	GpgContents string `json:"GpgContents"`
}

var req = Request{}

// ParseRequest get the json message body and map it to our Request struct.
// This should be the first middleware step.
func ParseRequest(w http.ResponseWriter, rc http.Request) (error, http.Request) {

	json := json.NewDecoder(rc.Body)

	err := json.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		return errors.New("Invalid JSON request in body"), rc
	}
	log.Println(req.Token)
	context := context.WithValue(rc.Context(), "request", req)
	return nil, *rc.WithContext(context)
}

