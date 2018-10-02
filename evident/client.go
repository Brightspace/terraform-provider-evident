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
	"net/http"
	"time"
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
	Credentials Credentials
}

type ExternalAccount struct {
	ID         string                    `json:"id"`
	Attributes ExternalAccountAttributes `json:"attributes"`
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
	utc := now.Format(time.RFC1123)
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

func makeRequest(request EvidentRequest, creds Credentials) (string, error) {
	baseURL := "https://api.evident.io"
	client := &http.Client{}

	req, err := http.NewRequest(request.Method, baseURL+request.URL, bytes.NewBuffer(request.Contents))
	if err != nil {
		return "", fmt.Errorf("Error creating request: %s", err)
	}

	location, _ := time.LoadLocation("GMT")
	t := time.Now().In(location)

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

	contentType := resp.Header.Get("Content-Type")
	if contentType == "application/vnd.api+json" {
		return "", fmt.Errorf("Content-Type is not a json type. Got: %s", contentType)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error while reading response body. %s", err)
	}

	return string(bytes), nil
}

func (evident Evident) all() ([]ExternalAccount, error) {
	var response EvidentResponse
	var result []ExternalAccount

	request := EvidentRequest{
		Method:   "GET",
		URL:      "/api/v2/external_accounts",
		Contents: []byte(""),
	}

	resp, err := makeRequest(request, evident.Credentials)
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

func (evident Evident) get(account string) (ExternalAccount, error) {
	var response EvidentResponse
	var result ExternalAccount

	request := EvidentRequest{
		Method:   "GET",
		URL:      "/api/v2/external_accounts/" + account,
		Contents: []byte(""),
	}

	resp, err := makeRequest(request, evident.Credentials)
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

func (evident Evident) delete(account string) (bool, error) {
	request := EvidentRequest{
		Method:   "DELETE",
		URL:      "/api/v2/external_accounts/" + account,
		Contents: []byte(""),
	}

	_, err := makeRequest(request, evident.Credentials)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (evident Evident) add(name string, arn string, externalID string, teamID string) (ExternalAccount, error) {
	var response EvidentResponse
	var result ExternalAccount

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

	resp, err := makeRequest(request, evident.Credentials)
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