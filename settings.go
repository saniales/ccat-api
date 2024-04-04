package ccatapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

type settingsClient struct {
	*ClientConfig
}

// newSettingsClient creates a new client with the provided Options.
func newSettingsClient(config ClientConfig) *settingsClient {
	c := &settingsClient{
		ClientConfig: &config,
	}

	WithBaseURL(fmt.Sprintf("%s/%s", c.baseURL, "settings"))(c.ClientConfig)

	return c
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
	Name     string `json:"name"`
	Value    any    `json:"value"`
	Category string `json:"category"`
}

// CreateSettingResponse contains the response of the CreateSetting method.
type CreateSettingResponse struct {
	createSettingRawResponse
	UpdatedAt time.Time `json:"updated_at"`
}

// createSettingRawResponse contains the raw response of the CreateSetting method.
//
// The raw response has a unix timestamp to represent time.
type createSettingRawResponse struct {
	Name      string `json:"name"`
	Value     any    `json:"value"`
	Category  string `json:"category"`
	SettingID string `json:"setting_id"`
	UpdatedAt int64  `json:"updated_at"`
}

// CreateSetting creates a new setting in the database.
func (c *settingsClient) CreateSetting(payload CreateSettingPayload) (*CreateSettingResponse, error) {
	rawResponse, err := doRequest[CreateSettingPayload, createSettingRawResponse](*c.ClientConfig, http.MethodPost, "", nil, &payload)
	if err != nil {
		return nil, err
	}

	response := &CreateSettingResponse{
		createSettingRawResponse: *rawResponse,
		UpdatedAt:                time.Unix(rawResponse.UpdatedAt, 0),
	}

	return response, nil
}

// UpdateSettingPayload contains the payload for the UpdateSetting method.
type UpdateSettingPayload struct {
	Name     *string `json:"name,omitempty"`
	Value    any     `json:"value,omitempty"`
	Category *string `json:"category,omitempty"`
}

// updateSettingRawResponse contains the raw response of the UpdateSetting method.
//
// The raw response has a unix timestamp to represent time.
type updateSettingRawResponse struct {
	Name      string `json:"name"`
	Value     any    `json:"value"`
	Category  string `json:"category"`
	SettingID string `json:"setting_id"`
	UpdatedAt int64  `json:"updated_at"`
}

// UpdateSettingResponse contains the response of the UpdateSetting method.
type UpdateSettingResponse struct {
	updateSettingRawResponse
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateSetting updates a specific setting in the database if it exists.
func (c *settingsClient) UpdateSetting(settingID string, payload UpdateSettingPayload) (*UpdateSettingResponse, error) {
	pathParams := fmt.Sprintf("/%s", settingID)
	rawResponse, err := doRequest[UpdateSettingPayload, updateSettingRawResponse](*c.ClientConfig, http.MethodPut, pathParams, nil, &payload)
	if err != nil {
		return nil, err
	}

	response := &UpdateSettingResponse{
		updateSettingRawResponse: *rawResponse,
		UpdatedAt:                time.Unix(rawResponse.UpdatedAt, 0),
	}

	return response, nil
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
