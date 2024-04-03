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

// Option is a functional option used to initialize a Client.
type Option func(*ClientConfig)

func WithConfig(newConfig ClientConfig) Option {
	return func(config *ClientConfig) {
		*config = newConfig
	}
}

// WithHTTPClient returns an Option function that sets the HTTP client for the Client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(config *ClientConfig) {
		if httpClient == nil {
			config.httpClient = http.DefaultClient
		} else {
			config.httpClient = httpClient
		}
	}
}

// WithBaseURL returns an Option function that sets the base URL of a Client.
func WithBaseURL(baseURL string) Option {
	return func(config *ClientConfig) {
		if (baseURL == "") || (baseURL == "/") {
			config.baseURL = defaultURL
		} else {
			config.baseURL = baseURL
		}
	}
}

// WithUserAgent returns an Option function that sets the user agent for the Client.
func WithUserAgent(userAgent string) Option {
	return func(config *ClientConfig) {
		if userAgent == "" {
			config.userAgent = defaultUserAgent
		}

		config.userAgent = userAgent
	}
}

// WithUserID returns an Option function that sets the user ID for the Client.
func WithUserID(userID string) Option {
	return func(config *ClientConfig) {
		config.userID = userID
	}
}

// WithAuthKey returns an Option function that sets the auth key for the Client.
func WithAuthKey(authKey string) Option {
	return func(config *ClientConfig) {
		config.authKey = authKey
	}
}

// WithMarshalFunc returns an Option function that sets the marshal function for the Client.
//
// if marshalFunc is nil, its default value is json.Marshal.
func WithMarshalFunc(marshalFunc func(v any) ([]byte, error)) Option {
	return func(config *ClientConfig) {
		if marshalFunc == nil {
			config.marshalFunc = json.Marshal
		} else {
			config.marshalFunc = marshalFunc
		}
	}
}

// WithUnmarshalFunc returns an Option function that sets the unmarshal function for the Client.
//
// if unmarshalFunc is nil, its default value is json.Unmarshal.
func WithUnmarshalFunc(unmarshalFunc func(data []byte, v any) error) Option {
	return func(config *ClientConfig) {
		if unmarshalFunc == nil {
			config.unmarshalFunc = json.Unmarshal
		} else {
			config.unmarshalFunc = unmarshalFunc
		}
	}
}
