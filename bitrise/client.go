package bitrise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

const baseBitriseURL = "https://api.bitrise.io/v0.1"

type Client struct {
	Token string
}

func (client *Client) ListBuilds(appSlug string, workflow string, branch string) ([]Build, error) {
	params := map[string]string{
		"workflow": workflow,
		"branch":   branch,
	}

	path := fmt.Sprintf("apps/%s/builds", appSlug)
	var buildsResponse BuildResponse
	err := client.apiRequest(http.MethodGet, path, params, nil, &buildsResponse)
	if err != nil {
		return nil, err
	}

	return buildsResponse.Data, nil
}

func (client *Client) apiRequest(method string, path string, params map[string]string, body map[string]interface{}, jsonResponse interface{}) error {
	headers := make(http.Header)
	headers.Add("Authorization", client.Token)

	endpointURL := prepareURL(path, params)
	bodyReader, err := prepareBody(method, &headers, body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(method, endpointURL.String(), bodyReader)
	if err != nil {
		return err
	}

	request.Header = headers

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	} else if !(response.StatusCode >= 200 && response.StatusCode < 300) {
		return fmt.Errorf("client: invalid status code %d", response.StatusCode)
	}

	json.NewDecoder(response.Body).Decode(&jsonResponse)

	return nil
}

func prepareBody(method string, headers *http.Header, body map[string]interface{}) (io.Reader, error) {
	if method == http.MethodGet || method == http.MethodDelete {
		return nil, nil
	} else if method == http.MethodPost || method == http.MethodPatch {
		headers.Add("Content-Type", "application/json")
		encodedBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		return bytes.NewBuffer(encodedBody), nil
	} else {
		return nil, fmt.Errorf("client: Unsupported http method %s", method)
	}
}

func prepareURL(urlPath string, params map[string]string) *url.URL {
	u, err := url.Parse(baseBitriseURL)
	if err != nil {
		log.Fatal(err)
	}

	u.Path = path.Join(u.Path, urlPath)
	values := u.Query()
	for k, v := range params {
		values.Add(k, v)
	}
	u.RawQuery = values.Encode()

	return u
}
