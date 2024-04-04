package ccatapi

import (
	"fmt"
	"net/http"
)

// llmClient is a sub-client for the LLM API.
type llmClient struct {
	Config *ClientConfig
}

// newLLMClient creates a new LLM sub-client with the provided config.
func newLLMClient(config ClientConfig) *llmClient {
	client := &llmClient{
		Config: &config,
	}

	WithBaseURL(fmt.Sprintf("%s/%s", client.Config.baseURL, "llm"))(client.Config)

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
	Description       string                              `json:"description"`
	HumanReadableName string                              `json:"humanReadableName"`
	Link              string                              `json:"link"`
	Properties        map[string]LLMSettingSchemaProperty `json:"properties"`
	Required          []string                            `json:"required"`
	Title             string                              `json:"title"`
	Type              string                              `json:"type"`
	LanguageModelName string                              `json:"languageModelName"`
}

// LLMSettingSchemaProperty contains the data about a single LLM setting schema property.
type LLMSettingSchemaProperty struct {
	Title   string `json:"title"`
	Type    string `json:"type"`
	Default *any   `json:"default,omitempty"`
}

// GetAllLLMsSettings returns a list of all LLMs settings.
func (client *llmClient) GetAllLLMsSettings() (*GetAllLLMsSettingsResponse, error) {
	resp, err := doRequest[any, GetAllLLMsSettingsResponse](*client.Config, http.MethodGet, "/settings", nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetLLMSetting returns a specific LLM setting.
func (client *llmClient) GetLLMSetting(languageModelName string) (*LLMSetting, error) {
	pathParams := fmt.Sprintf("/settings/%s", languageModelName)
	resp, err := doRequest[any, LLMSetting](*client.Config, http.MethodGet, pathParams, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UpsertLLMSetting updates a specific LLM setting value.
func (client *llmClient) UpsertLLMSetting(languageModelName string, value map[string]any) (*LLMSetting, error) {
	pathParams := fmt.Sprintf("/settings/%s", languageModelName)
	resp, err := doRequest[map[string]any, LLMSetting](*client.Config, http.MethodPut, pathParams, nil, &value)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
