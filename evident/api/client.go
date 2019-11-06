package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type Credentials struct {
	AccessKey []byte
	SecretKey []byte
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

type allExternalAccountAws struct {
	Data []ExternalAccount `json:"data"`
}

type ExternalAccount struct {
	ID         interface{} `json:"id"`
	Attributes struct {
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

		// Retry
		rest.SetRetryCount(evident.RetryMaximum)
		rest.SetRetryWaitTime(60 * time.Second)
		rest.SetRetryMaxWaitTime((5 * 60) * time.Second)
		rest.AddRetryCondition(func(r *resty.Response, err error) bool {
			switch code := r.StatusCode(); code {
			case http.StatusTooManyRequests:
				return true
			default:
				return false
			}
		})
		
		// Error handling
		rest.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
			status := r.StatusCode()
			if (status < 200) || (status >= 400) {
				return fmt.Errorf("Response not successful: Received status code %d.", status)
			}

			return nil
		})

		//Authentication
		rest.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			t := time.Now().UTC()
			key := string(evident.Credentials.AccessKey)
			secret := string(evident.Credentials.SecretKey)
			sign, _ := NewHTTPSignature(r.URL, r.Method, []byte(r.Body.(string)), t, key, secret)
			for name, value := range sign {
				r.SetHeader(name, value.(string))
			}
			return nil
		})
		evident.RestClient = rest
	}
	return evident.RestClient
}

func (evident *Evident) All() ([]ExternalAccount, error) {
	var result []ExternalAccount
	restClient := evident.GetRestClient()

	url := "/api/v2/external_accounts"
	req := restClient.R().SetBody("").SetResult(&allExternalAccountAws{})
	resp, err := req.Get(url)
	if err != nil {
		return result, err
	}

	response := resp.Result().(*allExternalAccountAws)
	if response != nil {
		result = response.Data
	}

	return result, nil
}

func (evident *Evident) Get(account string) (*ExternalAccount, error) {
	restClient := evident.GetRestClient()

	url := fmt.Sprintf("/api/v2/external_accounts/%s", account)
	req := restClient.R().SetBody("").SetResult(&getExternalAccountAws{})

	resp, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	response := resp.Result().(*getExternalAccountAws)
	if response == nil {
		return nil, nil
	}

	return &response.Data, nil
}

func (evident *Evident) Delete(account string) (bool, error) {
	restClient := evident.GetRestClient()

	url := fmt.Sprintf("/api/v2/external_accounts/%s", account)
	req := restClient.R().SetBody("")

	_, err := req.Delete(url)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (evident *Evident) Add(name string, arn string, externalID string, teamID string) (ExternalAccount, error) {
	var result ExternalAccount
	restClient := evident.GetRestClient()

	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "external_accounts",
			"attributes": map[string]string{
				"name":        name,
				"arn":         arn,
				"team_id":     teamID,
				"external_id": externalID,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return result, err
	}

	url := "/api/v2/external_accounts"
	req := restClient.R().SetBody(string(body)).SetResult(&getExternalAccountAws{})

	resp, err := req.Post(url)
	if err != nil {
		return result, err
	}

	response := resp.Result().(*getExternalAccountAws)
	response.Data.Attributes.Arn = arn
	response.Data.Attributes.ExternalID = externalID
	return response.Data, nil
}

func (evident *Evident) Update(account string, name string, arn string, externalID string, teamID string) (ExternalAccount, error) {
	var result ExternalAccount
	restClient := evident.GetRestClient()

	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "external_accounts",
			"attributes": map[string]string{
				"name":        name,
				"arn":         arn,
				"team_id":     teamID,
				"external_id": externalID,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return result, err
	}

	url := fmt.Sprintf("/api/v2/external_accounts/%+v", account)
	req := restClient.R().SetBody(string(body)).SetResult(&getExternalAccountAws{})

	resp, err := req.Patch(url)
	if err != nil {
		return result, err
	}

	response := resp.Result().(*getExternalAccountAws)
	return response.Data, nil
}
