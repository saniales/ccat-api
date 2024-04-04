package ccatapi

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// settingsClient is a sub-client for the Settings API.
type settingsClient struct {
	*ClientConfig
}

// newSettingsClient creates a new Settings sub-client with the provided config.
func newSettingsClient(config ClientConfig) *settingsClient {
	client := &settingsClient{
		ClientConfig: &config,
	}

	WithBaseURL(fmt.Sprintf("%s/%s", client.ClientConfig.baseURL, "settings"))(client.ClientConfig)

	return client
}

// GetSettingsParams contains the parameters for the GetSettings method.
type GetSettingsParams struct {
	Search string `url:"search,omitempty"`
}

// SettingsResponse contains the response of the GetSettings method.
type SettingsResponse struct {
	Settings []Setting `json:"settings"`
}

// SettingResponse contains the data about a single setting.
type Setting struct {
	Name      string         `json:"name"`
	Value     map[string]any `json:"value"`
	Category  string         `json:"category"`
	SettingID string         `json:"setting_id"`
	UpdatedAt int64          `json:"updated_at"`
}

// GetSettings returns a list of settings, optionally filtered by a search query.
func (c *settingsClient) GetSettings(params GetSettingsParams) (*SettingsResponse, error) {
	values, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	resp, err := doRequest[any, SettingsResponse](*c.ClientConfig, http.MethodGet, "", values, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreateSettingPayload contains the payload for the CreateSetting method.
type CreateSettingPayload struct {
	Name     string         `json:"name"`
	Value    map[string]any `json:"value"`
	Category string         `json:"category"`
}

// CreateSettingResponse contains the response of the CreateSetting method.
type CreateSettingResponse Setting

// CreateSetting creates a new setting in the database.
func (c *settingsClient) CreateSetting(payload CreateSettingPayload) (*CreateSettingResponse, error) {
	resp, err := doRequest[CreateSettingPayload, CreateSettingResponse](*c.ClientConfig, http.MethodPost, "", nil, &payload)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UpdateSettingPayload contains the payload for the UpdateSetting method.
type UpdateSettingPayload struct {
	Name     *string `json:"name,omitempty"`
	Value    any     `json:"value,omitempty"`
	Category *string `json:"category,omitempty"`
}

// UpdateSettingResponse contains the response of the UpdateSetting method.
type UpdateSettingResponse Setting

// UpdateSetting updates a specific setting in the database if it exists.
func (c *settingsClient) UpdateSetting(settingID string, payload UpdateSettingPayload) (*UpdateSettingResponse, error) {
	pathParams := fmt.Sprintf("/%s", settingID)
	resp, err := doRequest[UpdateSettingPayload, UpdateSettingResponse](*c.ClientConfig, http.MethodPut, pathParams, nil, &payload)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DeleteSEtting deletes a specific setting in the database.
func (c *settingsClient) DeleteSetting(settingID string) error {
	pathParams := fmt.Sprintf("/%s", settingID)
	_, err := doRequest[any, any](*c.ClientConfig, http.MethodDelete, pathParams, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
