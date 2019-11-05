package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
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
	RestClient   *resty.Client
	Credentials  Credentials
	RetryMaximum int
}

type getExternalAccountAws struct {
	Data ExternalAccount `json:"data"`
}

type ExternalAccount struct {
	ID         interface{} `json:"id"`
	Attributes  struct {
		Name       string `json:"name"`
		Provider   string `json:"provider"`
		Arn        string `json:"arn"`
		Account    string `json:"account"`
		ExternalID string `json:"external_id"`
	} `json:"attributes"`
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
	headers, _ := NewHTTPSignature(request.URL, request.Method, request.Contents, t, string(creds.AccessKey), string(creds.SecretKey))
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

func (evident *Evident) GetRestClient() *resty.Client {
	if evident.RestClient == nil {
		rest := resty.New()
		rest.SetHostURL("https://api.evident.io")

		evident.RestClient = rest
	}
	return evident.RestClient
}

func (evident *Evident) All() ([]ExternalAccount, error) {
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

func (evident *Evident) Get(account string) (ExternalAccount, error) {
	var result ExternalAccount
	restClient := evident.GetRestClient()
	credentials := evident.Credentials

	url := fmt.Sprintf("/api/v2/external_accounts/%s", account)
	req := restClient.R().SetBody("").SetResult(&getExternalAccountAws{})
	sign, _ := NewHTTPSignature(url, "GET", []byte(""), time.Now().UTC(), string(credentials.AccessKey), string(credentials.SecretKey))
	for name, value := range sign {
		req = req.SetHeader(name, value.(string))
	}

	resp, err := req.Get(url)
	if err != nil {
		return result, err
	}
	
	response := resp.Result().(*getExternalAccountAws)
	return response.Data, nil
}

func (evident *Evident) Delete(account string) (bool, error) {
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

func (evident *Evident) Add(name string, arn string, externalID string, teamID string) (ExternalAccount, error) {
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

func (evident *Evident) Update(account string, name string, arn string, externalID string, teamID string) (ExternalAccount, error) {
	var response EvidentResponse
	var result ExternalAccount
	var err error
	var resp string

	// Update Payload is the same as create payload so we use the same
	// payload for requesting
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
