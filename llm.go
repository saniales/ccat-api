package ccatapi

import (
	"fmt"
	"net/http"
)

// llmClient is a sub-client for the LLM API.
type llmClient struct {
	config *ClientConfig
}

// newLLMClient creates a new LLM sub-client with the provided config.
func newLLMClient(config ClientConfig) *llmClient {
	client := &llmClient{
		config: &config,
	}

	WithBaseURL(fmt.Sprintf("%s/%s", client.config.baseURL, "llm"))(client.config)

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
func (client *llmClient) GetAllLLMsSettings() (*GetAllLLMsSettingsResponse, error) {
	resp, err := doRequest[any, GetAllLLMsSettingsResponse](*client.config, http.MethodGet, "/settings", nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetLLMSetting returns a specific LLM setting.
func (client *llmClient) GetLLMSetting(languageModelName string) (*LLMSetting, error) {
	pathParams := fmt.Sprintf("/settings/%s", languageModelName)
	resp, err := doRequest[any, LLMSetting](*client.config, http.MethodGet, pathParams, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UpsertLLMSetting updates a specific LLM setting value.
func (client *llmClient) UpsertLLMSetting(languageModelName string, value map[string]any) (*LLMSetting, error) {
	pathParams := fmt.Sprintf("/settings/%s", languageModelName)
	resp, err := doRequest[map[string]any, LLMSetting](*client.config, http.MethodPut, pathParams, nil, &value)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
