// request package will parse all requests
package request

import (
	"encoding/json"
	"net/http"
	"errors"
	"context"
)

// Request put fields for all requests in this struct.
type Request struct {
	Token string `json:"token"`
	Scope []string `json:"Scope"`
	GpgContents string `json:"GpgContents"`
	Key string `json:"key"`
        Name    string `json:"name"`
        Comment string `json:"comment"`
        Email   string `json:"email"`
}

var req = Request{}

// ParseRequest get the json message body and map it to our Request struct.
// This is the only way to get information from the request body through the middlware chain.  The message body is destroyed after the first read here.
func ParseRequest(w http.ResponseWriter, rc http.Request) (error, http.Request) {

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

