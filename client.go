package ccatapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	*ClientConfig
}

type ClientConfig struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string

	userID  string
	authKey string

	marshalFunc   func(v any) ([]byte, error)
	unmarshalFunc func(data []byte, v any) error

	Settings *settingsClient
}

// NewClient creates a new client with the provided Options.
func NewClient(opts ...Option) *Client {
	c := &Client{
		ClientConfig: &ClientConfig{
			httpClient: http.DefaultClient,
			baseURL:    defaultURL,
			userAgent:  defaultUserAgent,

			userID:  defaultUserID,
			authKey: "",

			marshalFunc:   json.Marshal,
			unmarshalFunc: json.Unmarshal,
		},
	}

	for _, opt := range opts {
		opt(c.ClientConfig)
	}

	return c
}

func (c *Client) Status() error {
	_, err := doRequest[any, any](c, http.MethodGet, "", nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func doRequest[PayloadType any, ResponseType any](config ClientConfig, method string, path string, queryParams url.Values, payload *PayloadType) (*ResponseType, error) {
	fullURL, err := url.Parse(fmt.Sprintf("%s/%s", config.baseURL, path))
	if err != nil {
		return nil, err
	}

	if queryParams != nil {
		fullURL.RawQuery = queryParams.Encode()
	}

	var requestBodyBuffer *bytes.Buffer
	if payload != nil {
		encodedPayload, err := config.marshalFunc(payload)
		if err != nil {
			return nil, err
		}

		requestBodyBuffer = bytes.NewBuffer(encodedPayload)
	}

	req, err := http.NewRequest(method, fullURL.String(), requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", config.userAgent)

	if len(config.authKey) > 0 {
		req.Header.Set("Authorization", config.authKey)
	}

	resp, err := config.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := new(ResponseType)
	err = c.unmarshalFunc(respBodyBytes, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
