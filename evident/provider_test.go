package evident

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testEvidentProviders map[string]terraform.ResourceProvider
var testEvidentProvider *schema.Provider
var fakeHttpServer *httptest.Server
var fakeTestHandler http.HandlerFunc

var fakeId string = "1232193123"
var fakeArn string = "fakearn"
var fakeName string = "fakename"
var fakeExternalId string = "676767666"
var fakeTeamId string = "1231255543"

var updatedFakeArn string = "updatedFakearn"
var updatedFakeExternalId string = "12345678966"
var updatedFakeTeamId string = "443432324234"

var state string = ""

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func init() {
	fakeHttpServer = httptest.NewServer(fakeTestHandler)
	testEvidentProvider = Provider().(*schema.Provider)
	testEvidentProvider.ConfigureFunc = testConfigureFunction
	testEvidentProviders = map[string]terraform.ResourceProvider{
		"evident": testEvidentProvider,
	}
}
func updateState(newState string) {
	fmt.Printf("\nupdating state from %+v -> %+v\n", state, newState)
	state = newState
}
func testConfigureFunction(d *schema.ResourceData) (interface{}, error) {
	status := 200
	bodyString := ""

	httpClient := NewTestClient(func(r *http.Request) *http.Response {
		if r.Method == "POST" {
			updateState("created")
			bodyString = GetTestOkResponse(fakeId, fakeArn, fakeExternalId, fakeTeamId)
			status = 200
		} else if r.Method == "PATCH" {
			updateState("updated")
			bodyString = GetTestOkResponse(fakeId, updatedFakeArn, updatedFakeExternalId, updatedFakeTeamId)
		} else if r.Method == "GET" {
			if state == "created" {
				bodyString = GetTestOkResponse(fakeId, fakeArn, fakeExternalId, fakeTeamId)
			} else if state == "updated" {
				bodyString = GetTestOkResponse(fakeId, updatedFakeArn, updatedFakeExternalId, updatedFakeTeamId)
			} else {
				bodyString = "{}"
				status = 404
			}
		} else if r.Method == "DELETE" {
			updateState("deleted")
			bodyString = "{}"
		}
		return &http.Response{
			StatusCode: status,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(bodyString)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	client := Evident{
		Credentials: Credentials{
			AccessKey: []byte(d.Get("access_key").(string)),
			SecretKey: []byte(d.Get("secret_key").(string)),
		},
		RetryMaximum: 1,
	}
	client.SetHttpClient(httpClient)
	config := Config{
		EvidentClient: client,
	}

	return &config, nil

}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

var testProviders = map[string]terraform.ResourceProvider{
	"evident": Provider(),
}
