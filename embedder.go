package ccatapi

import (
	"fmt"
	"net/http"
)

// embedderClient is a sub-client for the Embedder API.
type embedderClient struct {
	config ClientConfig
}

// newEmbedderClient creates a new Embedder sub-client with the provided config.
func newEmbedderClient(config ClientConfig) *embedderClient {
	client := &embedderClient{
		config: config,
	}

	WithBaseURL(fmt.Sprintf("%s/%s", client.config.baseURL, "embedder"))(&client.config)

	return client
}

// GetAllEmbeddersSettingsResponse contains the response of the GetAllEmbeddersSettings method.
type GetAllEmbeddersSettingsResponse struct {
	Settings []EmbedderSetting `json:"settings"`
}

// EmbedderSetting contains the data about a single embedder setting.
type EmbedderSetting struct {
	Name   string                `json:"name"`
	Value  map[string]any        `json:"value"`
	Schema EmbedderSettingSchema `json:"schema"`
}

// EmbedderSettingSchema contains the JSON schema about a single embedder setting.
type EmbedderSettingSchema struct {
	settingSchema
	LanguageEmbedderName string `json:"languageEmbedderName"`
}

// GetAllEmbeddersSettings returns a list of all embedders settings.
func (client *embedderClient) GetAllEmbeddersSettings() (*GetAllEmbeddersSettingsResponse, error) {
	resp, err := doRequest[any, GetAllEmbeddersSettingsResponse](client.config, http.MethodGet, "/settings", nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
