package ccatapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

	c.Settings = newSettingsClient(*c.ClientConfig)

	return c
}

// Status returns the status of the Ccat API.
func (c *Client) Status() error {
	_, err := doRequest[any, any](*c.ClientConfig, http.MethodGet, "", nil, nil)
	if err != nil {
		return err
	}

	return nil
}

type APIErrorResponse struct {
	Errors []APIError `json:"error"`
}

func (err *APIErrorResponse) Error() string {
	var builder strings.Builder

	builder.WriteString("API error: \n")
	for _, err := range err.Errors {
		builder.WriteString(err.Error())
	}

	return builder.String()
}

type APIError struct {
	Type     string            `json:"type"`
	Location []string          `json:"loc"`
	Message  string            `json:"msg"`
	Input    map[string]string `json:"input"`
	URL      string            `json:"url"`
}

func (err *APIError) Error() string {
	var builder strings.Builder

	builder.WriteString("type: ")
	builder.WriteString(err.Type)
	builder.WriteString("\nlocation: ")
	for _, loc := range err.Location {
		builder.WriteString(loc)
		builder.WriteString(" ")
	}

	builder.WriteString("\nmessage: ")
	builder.WriteString(err.Message)
	for fieldName, fieldType := range err.Input {
		builder.WriteString("\n")
		builder.WriteString(fieldName)
		builder.WriteString(": ")
		builder.WriteString(fieldType)
	}

	builder.WriteString("\nurl: ")
	builder.WriteString(err.URL)

	return builder.String()
}

// doRequest sends a generic request to the Ccat API and returns the response.
//
// It uses a client config to keep consistency between clients.
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

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var apiErr APIErrorResponse
		err = config.unmarshalFunc(respBodyBytes, &apiErr)
		if err != nil {
			return nil, ErrUnknownError(resp.StatusCode, string(respBodyBytes))
		}

		return nil, &apiErr
	}

	response := new(ResponseType)
	err = config.unmarshalFunc(respBodyBytes, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
