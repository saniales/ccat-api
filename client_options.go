package ccatapi

import (
	"encoding/json"
	"net/http"
)

const (
	defaultURL       string = "http://localhost:1865"
	defaultUserAgent string = "ccat-api"
	defaultUserID    string = "user"
)

// option is a functional option used to initialize a Client.
type option func(*clientConfig)

// WithHTTPClient returns an option function that sets the HTTP client for the Client.
func WithHTTPClient(httpClient *http.Client) option {
	return func(config *clientConfig) {
		if httpClient == nil {
			config.httpClient = http.DefaultClient
		} else {
			config.httpClient = httpClient
		}
	}
}

// WithBaseURL returns an option function that sets the base URL of a Client.
func WithBaseURL(baseURL string) option {
	return func(config *clientConfig) {
		if (baseURL == "") || (baseURL == "/") {
			config.baseURL = defaultURL
		} else {
			config.baseURL = baseURL
		}
	}
}

// WithUserAgent returns an option function that sets the user agent for the Client.
func WithUserAgent(userAgent string) option {
	return func(config *clientConfig) {
		if userAgent == "" {
			config.userAgent = defaultUserAgent
		}

		config.userAgent = userAgent
	}
}

// WithUserID returns an option function that sets the user ID for the Client.
func WithUserID(userID string) option {
	return func(config *clientConfig) {
		config.userID = userID
	}
}

// WithAuthKey returns an option function that sets the auth key for the Client.
func WithAuthKey(authKey string) option {
	return func(config *clientConfig) {
		config.authKey = authKey
	}
}

// WithMarshalFunc returns an option function that sets the marshal function for the Client.
//
// if marshalFunc is nil, its default value is json.Marshal.
func WithMarshalFunc(marshalFunc func(v any) ([]byte, error)) option {
	return func(config *clientConfig) {
		if marshalFunc == nil {
			config.marshalFunc = json.Marshal
		} else {
			config.marshalFunc = marshalFunc
		}
	}
}

// WithUnmarshalFunc returns an option function that sets the unmarshal function for the Client.
//
// if unmarshalFunc is nil, its default value is json.Unmarshal.
func WithUnmarshalFunc(unmarshalFunc func(data []byte, v any) error) option {
	return func(config *clientConfig) {
		if unmarshalFunc == nil {
			config.unmarshalFunc = json.Unmarshal
		} else {
			config.unmarshalFunc = unmarshalFunc
		}
	}
}
