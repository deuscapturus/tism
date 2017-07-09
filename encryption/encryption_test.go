// token Testing
package encryption

import (
	"context"
	"crypto"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// TestParse tests the Parse middleware function.  Verify http status codes and updated context with user input.
func TestSetMyKeyRing(t *testing.T) {

	// Variables stub/mock
	// Need to make fake keyring.
	cases := []struct {
		name    string
		comment string
		email   string
		claims  string
	}{
		{"test1", "comment1", "test1@tests.test", "limited"},
		//		{"test2", "comment2", "test2@tests.test", "all"},
	}

	// Create new OpenGPG entity
	pgpConfig := &packet.Config{
		DefaultHash: crypto.SHA256,
	}

	for _, c := range cases {
		NewEntity, err := openpgp.NewEntity(c.name, c.comment, c.email, pgpConfig)
		if err != nil {
			t.Fatal(err)
		}
		// Turn NewEntity into an entiylist and pass that into KeyRing
		KeyRing = MyEntityList{openpgp.EntityList{NewEntity}}
		NewEntityId := strconv.FormatUint(NewEntity.PrimaryKey.KeyId, 16)
		TestClaims := []string{NewEntityId}

		// Create a stub/mock request with http.NewRequest
		req, err := http.NewRequest("POST", "/", nil)
		if err != nil {
			t.Fatalf("%v: %v", c.name, err)
		}

		ctx := req.Context()

		if c.claims == "limited" {
			ctx = context.WithValue(ctx, "claims", TestClaims)
			ctx = context.WithValue(ctx, "claimsAll", false)
		} else if c.claims == "all" {
			ctx = context.WithValue(ctx, "claims", TestClaims)
			ctx = context.WithValue(ctx, "claimsAll", true)
		} else {
			t.Fatalf("claims must be limited or all")
		}

		req = req.WithContext(ctx)

		// Create a test response recorder
		res := httptest.NewRecorder()

		// Create http handler wrapper.
		// The middleware function returns an error and http.Request,
		// so we can't use it directly in http.HandlerFunc.
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Test Parse function for errors
			err, rc := SetMyKeyRing(w, *r)
			if err != nil {
				t.Fatalf("%v: %v", c.name, err)
			}

			// Test that keyring is found in context of rc
			FoundKeyRing := rc.Context().Value("MyKeyRing").(openpgp.EntityList)

			if len(FoundKeyRing) == 0 {
				t.Fatal("MyKeyRing is empty")
			} else if FoundKeyRing[0] != NewEntity {
				t.Error("Generated keyring does not match the keyring found in returned request context")
			}
		})

		handler.ServeHTTP(res, req)
	}
}
