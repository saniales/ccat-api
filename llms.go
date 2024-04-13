package ccatapi

import (
	"fmt"
	"net/http"
)

// llmsClient is a sub-client for the LLM API.
type llmsClient struct {
	config clientConfig
}

// newLLMsClient creates a new LLM sub-client with the provided config.
func newLLMsClient(config clientConfig) *llmsClient {
	client := &llmsClient{
		config: config,
	}

	WithBaseURL(fmt.Sprintf("%s/%s", client.config.baseURL, "llm"))(&client.config)

	return client
}

type GetAllLLMsSettingsResponse struct {
	Settings []LLMSetting `json:"settings"`
}

// LLMSetting contains the data about a single LLM setting.
type LLMSetting struct {
	Name   string           `json:"name"`
	Value  map[string]any   `json:"value"`
	Schema LLMSettingSchema `json:"schema"`
}

// LLMSettingSchema contains the data about a single LLM setting schema.
type LLMSettingSchema struct {
	settingSchema
	LanguageModelName string `json:"languageModelName"`
}

// GetAllLLMsSettings returns a list of all LLMs settings.
func (client *llmsClient) GetAllLLMsSettings() (*GetAllLLMsSettingsResponse, error) {
	resp, err := doAPIRequest[any, GetAllLLMsSettingsResponse](
		client.config,
		http.MethodGet,
		"/settings",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetLLMSetting returns a specific LLM setting.
func (client *llmsClient) GetLLMSetting(languageModelName string) (*LLMSetting, error) {
	pathParams := fmt.Sprintf("/settings/%s", languageModelName)

	resp, err := doAPIRequest[any, LLMSetting](
		client.config,
		http.MethodGet,
		pathParams,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UpsertLLMSetting updates a specific LLM setting value.
func (client *llmsClient) UpsertLLMSetting(languageModelName string, value map[string]any) (*LLMSetting, error) {
	pathParams := fmt.Sprintf("/settings/%s", languageModelName)

	resp, err := doAPIRequest[map[string]any, LLMSetting](
		client.config,
		http.MethodPut,
		pathParams,
		nil,
		&value,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
