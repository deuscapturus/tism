// token Testing
package token

import (
	"context"
	"github.com/deuscapturus/tism/config"
	"github.com/deuscapturus/tism/request"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

// TestNew tests the New middleware function.
func TestNew(t *testing.T) {

	// Variables stub/mock
	cases := []struct {
		admin      int
		keys       []string
		expiration int64
	}{
		{1, []string{"ALL"}, time.Now().Add(time.Hour * 30303).Unix()},
	}

	// Set mock settings
	config.Config.JWTsecret = "12345"

	for _, c := range cases {
		// Create a stub/mock request with http.NewRequest
		req, err := http.NewRequest("POST", "/", nil)
		if err != nil {
			t.Fatalf("%v", err)
		}
		reqContext := request.Request{Expiration: c.expiration, Admin: c.admin, Keys: c.keys}
		ctx := req.Context()
		ctx = context.WithValue(ctx, "request", reqContext)
		ctx = context.WithValue(ctx, "admin", 1)
		ctx = context.WithValue(ctx, "claims", "ALL")
		req = req.WithContext(ctx)

		// Create a test response recorder
		res := httptest.NewRecorder()

		// Create http handler wrapper.
		// The middleware function returns an error and http.Request,
		// so we can't use it directly in http.HandlerFunc.
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Test Parse function for errors
			err, _ := New(w, *r)
			if err != nil {
				t.Fatalf("%v", err)
			}
		})

		handler.ServeHTTP(res, req)
		//Confirm my jwt values match my case
		jwtBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("%v", err)
		}
		token, err := parseToken(string(jwtBytes))
		if err != nil {
			t.Fatalf("%v", err)
		}
		if token.Claims.(*JwtClaimsMap).Admin != c.admin {
			t.Errorf("Generated token admin value is not correct.  Expected: %v, Found: %v", c.admin, token.Claims.(*JwtClaimsMap).Admin)
		}
		if reflect.DeepEqual(token.Claims.(*JwtClaimsMap).Keys, c.keys) != true {
			t.Errorf("Generated token keys value is not correct.  Expected: %v, Found: %v", c.keys, token.Claims.(*JwtClaimsMap).Keys)
		}
		if token.Claims.(*JwtClaimsMap).Expiration != c.expiration {
			t.Errorf("Generated token expiration value is not correct.  Expected: %v, Found: %v", c.expiration, token.Claims.(*JwtClaimsMap).Expiration)
		}
	}
}

// TestParse tests the Parse middleware function.  Verify http status codes and updated context with user input.
func TestParse(t *testing.T) {

	// Variables stub/mock
	cases := []struct {
		name       string
		token      string
		claimstype reflect.Kind
	}{
		{"ALL-Admin", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjo5OTk5OTk5OTk5OSwianRpIjoiYTA5ZDIydmw1dG83Iiwia2V5cyI6WyJBTEwiXX0.73T6TWAlNcv4Jt_HltjamLgHm0yF0M8XTUaWpgMLwy4", reflect.String},
		{"Limited-NonAdmin", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MCwiZXhwIjoxNjA0NjkyMDgwLCJqdGkiOiIxbm02MW9pbmEydTdnIiwia2V5cyI6WyI4MTVmOTlmOGY5ZDQzNWUzIiwiMTNlYzgwYzc1YzY5NzA1NSJdfQ.suObIX8YYVL0qCqfT_lmXDSSxTr8IsnXqKDxlnb8GXk", reflect.Slice},
		{"Limited-Admin", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNjA0NzIwOTM3LCJqdGkiOiJmamphZnF0b2hhaWsiLCJrZXlzIjpbIjgxNWY5OWY4ZjlkNDM1ZTMiLCIxM2VjODBjNzVjNjk3MDU1Il19.WplBDakhsMOp786_NlOmIzWT8-VmXZInJ9jne6qsI40", reflect.Slice},
	}

	// Set mock settings
	config.Config.JWTsecret = "12345"

	for _, c := range cases {
		reqContext := request.Request{Token: c.token}

		// Create a stub/mock request with http.NewRequest
		req, err := http.NewRequest("POST", "/", nil)
		if err != nil {
			t.Fatalf("%v: %v", c.name, err)
		}
		ctx := req.Context()
		ctx = context.WithValue(ctx, "request", reqContext)
		req = req.WithContext(ctx)

		// Create a test response recorder
		res := httptest.NewRecorder()

		// Create http handler wrapper.
		// The middleware function returns an error and http.Request,
		// so we can't use it directly in http.HandlerFunc.
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Test Parse function for errors
			err, rc := Parse(w, *r)
			if err != nil {
				t.Fatalf("%v: %v", c.name, err)
			}

			// Test context value against the testing table.
			Foundtype := reflect.TypeOf(rc.Context().Value("claims")).Kind()
			if Foundtype != c.claimstype {
				t.Errorf("%v: Claims context type incorrect.  Expected %v, Found %v", c.name, Foundtype, c.claimstype)
			}
		})

		handler.ServeHTTP(res, req)
	}
}

// BenchmarkParse performance check for Parse.
func BenchmarkParse(b *testing.B) {

	// Variables stub/mock
	onecase := struct {
		name       string
		token      string
		claimstype reflect.Kind
	}{
		"Limited-Admin",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNjA0NzIwOTM3LCJqdGkiOiJmamphZnF0b2hhaWsiLCJrZXlzIjpbIjgxNWY5OWY4ZjlkNDM1ZTMiLCIxM2VjODBjNzVjNjk3MDU1Il19.WplBDakhsMOp786_NlOmIzWT8-VmXZInJ9jne6qsI40",
		reflect.Slice,
	}

	// Set mock settings
	config.Config.JWTsecret = "12345"

	// Create a test response recorder
	res := httptest.NewRecorder()

	// Create http handler.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err, _ := Parse(w, *r)
		if err != nil {
			b.Fatal(err)
		}

	})

	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		b.Fatal(err)
	}

	reqContext := request.Request{Token: onecase.token}
	ctx := req.Context()
	ctx = context.WithValue(ctx, "request", reqContext)
	req = req.WithContext(ctx)

	b.ResetTimer()
	// Test Parse function for errors
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(res, req)
	}

}
