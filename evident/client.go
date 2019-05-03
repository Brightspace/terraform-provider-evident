package evident

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/matryer/try"
)

type Credentials struct {
	AccessKey []byte
	SecretKey []byte
}

type EvidentRequest struct {
	Method   string
	URL      string
	Contents []byte
}

type EvidentResponse struct {
	Data json.RawMessage `json:data`
}

type Evident struct {
	HttpClient   *http.Client
	Credentials  Credentials
	RetryMaximum int
}

type ExternalAccount struct {
	ID         interface{}               `json:"id"`
	Attributes ExternalAccountAttributes `json:"attributes"`
}

func (ec *ExternalAccount) GetIdString() string {

	switch v := ec.ID.(type) {
	case float64:
		return fmt.Sprintf("%.0f", v)
	default:
		return fmt.Sprintf("%+v", v)
	}
}

type CmdAddExternalAccount struct {
	Data CmdAddExternalAccountPayload `json:"data"`
}

type CmdAddExternalAccountPayload struct {
	Type       string                          `json:"type"`
	Attributes CmdAddExternalAccountAttributes `json:"attributes"`
}

type CmdAddExternalAccountAttributes struct {
	Name       string `json:"name"`
	ExternalID string `json:"external_id"`
	ARN        string `json:"arn"`
	TeamID     string `json:"team_id"`
}

//TODO: FIX Logic:
//We don't have Name and Provider in the attributes they are in related entities
//so our resource assessment is wrong (only ids and arns work which should suffice for now)
//at some point we should integrate the jsonapi client into this
//also we don't add team id to resources which might cause troubles in the future

type ExternalAccountAttributes struct {
	Name       string `json:"name"`
	Provider   string `json:"provider"`
	Arn        string `json:"arn"`
	Account    string `json:"account"`
	ExternalID string `json:"external_id"`
}

func makeHeaders(now time.Time, req EvidentRequest, creds Credentials) (map[string]interface{}, error) {
	headers := make(map[string]interface{})

	ctype := "application/vnd.api+json"
	md5sum := md5.Sum(req.Contents)
	md5 := base64.StdEncoding.EncodeToString(md5sum[:])
	utc := now.Format(http.TimeFormat)
	request := fmt.Sprintf("%s,%s,%s,%s,%s", req.Method, ctype, md5, req.URL, utc)
	encodedAuth := makeAuth([]byte(request), []byte(creds.SecretKey))

	headers["Accept"] = ctype
	headers["Content-Type"] = ctype
	headers["Content-MD5"] = md5
	headers["Date"] = utc

	headers["Authorization"] = fmt.Sprintf("APIAuth %s:%s", creds.AccessKey, encodedAuth)

	return headers, nil
}

func makeAuth(message []byte, key []byte) string {
	hash := hmac.New(sha1.New, key)
	hash.Write(message)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (evident *Evident) makeRequest(request EvidentRequest, creds Credentials) (string, error) {
	baseURL := "https://api.evident.io"
	client := evident.GetHttpClient()
	reqURL := baseURL + request.URL

	log.Printf("[DEBUG] sending request: (Request: %q, URL: %q)", request.URL, reqURL)
	req, err := http.NewRequest(request.Method, reqURL, bytes.NewBuffer(request.Contents))
	if err != nil {
		return "", fmt.Errorf("Error creating request: %s", err)
	}

	t := time.Now().UTC()
	headers, _ := makeHeaders(t, request, creds)
	for name, value := range headers {
		req.Header.Set(name, value.(string))
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error during making a request: %s", request.URL)

	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP request error. Response code: %d", resp.StatusCode)

	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error while reading response body. %s", err)
	}

	return string(bytes), nil
}
func (evident *Evident) SetHttpClient(client *http.Client) {
	evident.HttpClient = client
}

func (evident *Evident) GetHttpClient() *http.Client {
	if evident.HttpClient == nil {
		evident.HttpClient = &http.Client{}
	}
	return evident.HttpClient
}

func (evident *Evident) all() ([]ExternalAccount, error) {
	var response EvidentResponse
	var result []ExternalAccount

	request := EvidentRequest{
		Method:   "GET",
		URL:      "/api/v2/external_accounts",
		Contents: []byte(""),
	}

	resp, err := evident.makeRequest(request, evident.Credentials)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response.Data), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (evident *Evident) get(account string) (ExternalAccount, error) {
	var response EvidentResponse
	var result ExternalAccount
	var err error
	var resp string

	request := EvidentRequest{
		Method:   "GET",
		URL:      "/api/v2/external_accounts/" + account,
		Contents: []byte(""),
	}

	err = try.Do(func(ampt int) (bool, error) {
		var err error
		resp, err = evident.makeRequest(request, evident.Credentials)
		if err != nil {
			log.Printf("[DEBUG] retrying request: (Attempt: %d/%d, URL: %q)", ampt, evident.RetryMaximum, err)
			time.Sleep(30 * time.Second)
		}
		return ampt < evident.RetryMaximum, err
	})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response.Data), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (evident *Evident) delete(account string) (bool, error) {
	var err error

	request := EvidentRequest{
		Method:   "DELETE",
		URL:      "/api/v2/external_accounts/" + account,
		Contents: []byte(""),
	}

	err = try.Do(func(ampt int) (bool, error) {
		var err error
		_, err = evident.makeRequest(request, evident.Credentials)
		if err != nil {
			log.Printf("[DEBUG] retrying request: (Attempt: %d/%d, URL: %q)", ampt, evident.RetryMaximum, err)
			time.Sleep(30 * time.Second)
		}
		return ampt < evident.RetryMaximum, err
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (evident *Evident) add(name string, arn string, externalID string, teamID string) (ExternalAccount, error) {
	var response EvidentResponse
	var result ExternalAccount
	var err error
	var resp string

	cmd := CmdAddExternalAccount{
		Data: CmdAddExternalAccountPayload{
			Type: "external_accounts",
			Attributes: CmdAddExternalAccountAttributes{
				Name:       name,
				ExternalID: externalID,
				TeamID:     teamID,
				ARN:        arn,
			},
		},
	}

	payloadJSON, err := json.Marshal(cmd)
	if err != nil {
		return result, err
	}

	request := EvidentRequest{
		Method:   "POST",
		URL:      "/api/v2/external_accounts",
		Contents: payloadJSON,
	}

	err = try.Do(func(ampt int) (bool, error) {
		var err error
		resp, err = evident.makeRequest(request, evident.Credentials)
		if err != nil {
			log.Printf("[DEBUG] retrying request: (Attempt: %d/%d, URL: %q)", ampt, evident.RetryMaximum, err)
			time.Sleep(30 * time.Second)
		}
		return ampt < evident.RetryMaximum, err
	})

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response.Data), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (evident *Evident) update(account string, name string, arn string, externalID string, teamID string) (ExternalAccount, error) {
	var response EvidentResponse
	var result ExternalAccount
	var err error
	var resp string

	cmd := CmdAddExternalAccount{
		Data: CmdAddExternalAccountPayload{
			Type: "external_accounts",
			Attributes: CmdAddExternalAccountAttributes{
				Name:       name,
				ExternalID: externalID,
				TeamID:     teamID,
				ARN:        arn,
			},
		},
	}

	payloadJSON, err := json.Marshal(cmd)
	if err != nil {
		return result, err
	}

	request := EvidentRequest{
		Method:   "PATCH",
		URL:      fmt.Sprintf("/api/v2/external_accounts/%+v", account),
		Contents: payloadJSON,
	}

	err = try.Do(func(ampt int) (bool, error) {
		var err error
		resp, err = evident.makeRequest(request, evident.Credentials)
		if err != nil {
			log.Printf("[DEBUG] retrying request: (Attempt: %d/%d, URL: %q)", ampt, evident.RetryMaximum, err)
			time.Sleep(30 * time.Second)
		}
		return ampt < evident.RetryMaximum, err
	})

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response.Data), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
