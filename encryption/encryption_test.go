// token Testing
package encryption

import (
	"context"
	"crypto"
	"github.com/deuscapturus/tism/request"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var testClearText string
var testBase64Cipher string
var testArmorCipher string
var NewEntityId string

func init() {
	testClearText = "pmja#ttA*RYRKlET*@l6ppN^VaVu4MwUeNj1Pz*AjA2JJ7y!pcJEmRWpahH!hvp%hdoLVYNWz!%RH&sXaNU0a3RzAPR6dehVgfQrmJkmkug*Z3hwIEEj%4J#KJrIuBQ21yRGoZxuxnpCD3ZuDmhNg!Ax"
	NewEntityId = mockKeyRing("test", "test comment", "test@email.com")
}

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

func mockKeyRing(name string, comment string, email string) (NewEntityId string) {

	// Create new OpenGPG entity
	pgpConfig := &packet.Config{
		DefaultHash: crypto.SHA256,
	}

	NewEntity, err := openpgp.NewEntity(name, comment, email, pgpConfig)
	if err != nil {
		panic(err)
	}
	// Turn NewEntity into an entiylist and pass that into KeyRing
	KeyRing = MyEntityList{openpgp.EntityList{NewEntity}}
	return strconv.FormatUint(NewEntity.PrimaryKey.KeyId, 16)
}

// TestEncrypt tests the Encrypt middleware function.  Verify http status codes validate decrypted text.
func TestEncrypt(t *testing.T) {

	cases := []struct {
		cleartext string
		encoding  string
	}{
		{testClearText, "base64"},
		{testClearText, "armor"},
	}

	for _, c := range cases {

		// Create a stub/mock request with http.NewRequest
		req, err := http.NewRequest("POST", "/", nil)
		if err != nil {
			t.Fatalf("%v: %v", c.encoding, err)
		}

		ctx := req.Context()

		MyKeyRing := KeyRing.EntityList
		ctx = context.WithValue(ctx, "MyKeyRing", MyKeyRing)

		var request = request.Request{}
		request.Id = strconv.FormatUint(MyKeyRing[0].PrimaryKey.KeyId, 16)
		request.DecSecret = testClearText

		if c.encoding == "base64" {
			request.Encoding = "base64"

		} else if c.encoding == "armor" {
			request.Encoding = "armor"
		}

		ctx = context.WithValue(ctx, "request", request)

		req = req.WithContext(ctx)

		// Create a test response recorder
		res := httptest.NewRecorder()

		// Create http handler wrapper.
		// The middleware function returns an error and http.Request,
		// so we can't use it directly in http.HandlerFunc.
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Test Parse function for errors
			err, _ := Encrypt(w, *r)
			if err != nil {
				t.Fatalf("%v: %v", c.encoding, err)
			}

		})

		handler.ServeHTTP(res, req)

	}

}
