package ccatapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// Client is a Cheshire Cat API client.
type Client struct {
	config clientConfig

	Settings  *settingsClient
	LLMs      *llmsClient
	Embedders *embeddersClient
	Plugins   *pluginsClient
}

// clientConfig is the configuration for the Cheshire Cat API client.
type clientConfig struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string

	userID  string
	authKey string

	marshalFunc   func(v any) ([]byte, error)
	unmarshalFunc func(data []byte, v any) error
}

// NewClient creates a new client with the provided Options.
func NewClient(opts ...option) *Client {
	client := &Client{
		config: clientConfig{
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
		opt(&client.config)
	}

	client.Settings = newSettingsClient(client.config)
	client.LLMs = newLLMsClient(client.config)
	client.Embedders = newEmbeddersClient(client.config)
	client.Plugins = newPluginsClient(client.config)

	return client
}

// Status returns the status of the Cheshire Cat API.
func (client *Client) Status() error {
	_, err := doAPIRequest[any, any](client.config, http.MethodGet, "", nil, nil)
	if err != nil {
		return err
	}

	return nil
}

type APIErrorsResponse struct {
	Errors []APIError `json:"error"`
}

func (err APIErrorsResponse) Error() string {
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

type APIErrorText struct {
	ErrorMessage string `json:"error"`
}

func (err APIErrorText) Error() string {
	return err.ErrorMessage
}

func (err APIError) Error() string {
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

// doAPIRequest sends a generic request to the Cheshire Cat API and returns the response.
//
// It uses a client config to keep consistency between clients.
func doAPIRequest[PayloadType any, ResponseType any](
	config clientConfig,
	method string,
	path string,
	queryParams url.Values,
	payload *PayloadType,
) (*ResponseType, error) {
	var requestBodyBuffer *bytes.Buffer
	if payload != nil {
		encodedPayload, err := config.marshalFunc(payload)
		if err != nil {
			return nil, err
		}

		requestBodyBuffer = bytes.NewBuffer(encodedPayload)
	}

	return doHTTPRequest[ResponseType](
		config,
		"application/json",
		method,
		path,
		queryParams,
		requestBodyBuffer,
	)
}

// doMultipartRequest performs a multipart request.
func doMultipartRequest[ResponseType any](config clientConfig, method string, path string, queryParams url.Values, payloadFileName string, payloadReader io.Reader) (*ResponseType, error) {
	if payloadReader == nil {
		return nil, ErrUploadMissingFile
	}

	var requestBodyBuffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBodyBuffer)

	w, err := multipartWriter.CreateFormFile("file", payloadFileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(w, payloadReader)
	if err != nil {
		return nil, err
	}

	multipartWriter.Close()

	return doHTTPRequest[ResponseType](
		config,
		multipartWriter.FormDataContentType(),
		method,
		path,
		queryParams,
		&requestBodyBuffer,
	)
}

// doHTTPRequest performs a generic raw HTTP request and returns the parsed response.
func doHTTPRequest[ResponseType any](
	config clientConfig,
	contentType string,
	method string,
	path string,
	queryParams url.Values,
	body io.Reader,
) (*ResponseType, error) {
	fullURL, err := url.Parse(fmt.Sprintf("%s/%s", config.baseURL, path))
	if err != nil {
		return nil, err
	}

	if queryParams != nil {
		fullURL.RawQuery = queryParams.Encode()
	}

	req, err := http.NewRequest(method, fullURL.String(), body)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", contentType)
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
		// Here we try to parse all possible error response formats
		// if any of them match, we return the corresponding error
		var apiErr error = new(APIErrorsResponse)
		err = config.unmarshalFunc(respBodyBytes, apiErr)
		if err == nil {
			return nil, apiErr
		}

		apiErr = new(APIErrorText)
		err = config.unmarshalFunc(respBodyBytes, apiErr)
		if err == nil {
			return nil, apiErr
		}

		apiErr = new(APIError)
		err = config.unmarshalFunc(respBodyBytes, apiErr)
		if err == nil {
			return nil, apiErr
		}

		// if none matches, we return an unknown error
		return nil, errUnknownError(resp.StatusCode, string(respBodyBytes))
	}

	response := new(ResponseType)
	err = config.unmarshalFunc(respBodyBytes, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
