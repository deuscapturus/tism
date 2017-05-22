// request Testing
package request

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestParse tests the Parse middleware function.  Verify http status codes and updated context with user input.
func TestParse(t *testing.T) {

	// Variables stub/mock
	token := "12345.abcde.98765"
	data := `{"token": "` + token + `"}"`

	// Create a stub/mock request with http.NewRequest
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(data))
	if err != nil {
		t.Fatal(err)
	}

	// Create a test response recorder
	res := httptest.NewRecorder()

	// Create http handler.
	// The middleware function also passes back an error, so we can't use it directly as http.HandlerFunc.
	// This is a wrapper.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc := *r.WithContext(r.Context())

		// Test Parse function for errors
		err, rc := Parse(w, rc)
		if err != nil {
			t.Fatal(err)
		}

		// Test that the context was updated.
		ctx := rc.Context().Value("request").(Request)
		// Check for Token.  In the future we may want to verify more input values
		if ctx.Token != token {
			t.Errorf("Request context not populated with token. Found: %v", ctx.Token)
		}
	})

	handler.ServeHTTP(res, req)
}
