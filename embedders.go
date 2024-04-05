package ccatapi

import (
	"fmt"
	"net/http"
)

// embeddersClient is a sub-client for the Embedder API.
type embeddersClient struct {
	config clientConfig
}

// newEmbeddersClient creates a new Embedder sub-client with the provided config.
func newEmbeddersClient(config clientConfig) *embeddersClient {
	client := &embeddersClient{
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
func (client *embeddersClient) GetAllEmbeddersSettings() (*GetAllEmbeddersSettingsResponse, error) {
	resp, err := doRequest[any, GetAllEmbeddersSettingsResponse](client.config, http.MethodGet, "/settings", nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}