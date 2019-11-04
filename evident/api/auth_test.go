package api

import (
	"bytes"
	"net/http"
	"testing"
	"time"
)

func TestNewHMAC(t *testing.T) {
	expected := "pFgU9Klwd8J36NJDl2aH4eqlG/4="
	key := []byte("SECRET")
	message := []byte("Hello! How are you?")

	actual := NewHMAC(message, key)
	if actual != expected {
		t.Errorf("HMAC was incorrect, actual: [%s], expected: [%s]", actual, expected)
	}
}

func TestNewHTTPSignatureHasExpectedHeaders(t *testing.T) {
	then := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	public := []byte("public")
	secret := []byte("secret")
	contents := []byte("Hello World!")
	req, _ := http.NewRequest("GET", "api.evident.io", bytes.NewBuffer(contents))

	expected := 5
	actual, _ := NewHTTPSignature(req, contents, then, public, secret)
	if len(actual) != expected {
		t.Errorf("HTTPHeaders length is incorrect, actual: [%d], expected: [%d]", len(actual), expected)
	}
}

func TestNewHTTPSignatureEnsuresHeaders(t *testing.T) {
	then := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	public := []byte("public")
	secret := []byte("secret")
	contents := []byte("Hello World!")
	req, _ := http.NewRequest("GET", "api.evident.io", bytes.NewBuffer(contents))

	keys := []string{"Accept", "Content-Type", "Content-MD5", "Date", "Authorization"}
	actual, _ := NewHTTPSignature(req, contents, then, public, secret)
	for _, element := range keys {
		if _, ok := actual[element]; !ok {
			t.Errorf("HTTPHeaders was missing header, expected: [%s]", element)
		}
	}
}
