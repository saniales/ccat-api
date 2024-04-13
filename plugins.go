package ccatapi

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// pluginsClient is a sub-client for the Plugins API.
type pluginsClient struct {
	config clientConfig
}

// newPluginsClient creates a new Plugins sub-client with the provided config.
func newPluginsClient(config clientConfig) *pluginsClient {
	client := &pluginsClient{
		config: config,
	}

	WithBaseURL(fmt.Sprintf("%s/%s", client.config.baseURL, "plugins"))(&client.config)

	return client
}

// PluginsResponse contains the response data about a plugin.
type PluginsResponse struct {
	Filters   pluginResponseFilters `json:"filters"`
	Installed []InstalledPlugin     `json:"installed"`
	Registry  []RegistryPlugin      `json:"registry"`
}

type pluginResponseFilters struct {
	Query *string `json:"query,omitempty"`
}

type plugin struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	AuthorName  string `json:"author_name"`
	AuthorURL   string `json:"author_url"`
	PluginURL   string `json:"plugin_url"`
	Tags        string `json:"tags"`
	Thumb       string `json:"thumb"`
}

// InstalledPlugin contains the data about an installed plugin.
type InstalledPlugin struct {
	plugin
	ID      string                    `json:"id"`
	Active  bool                      `json:"active"`
	Upgrade *string                   `json:"upgrade"`
	Hooks   []InstalledPluginHookData `json:"hooks"`
	Tools   []InstalledPluginToolData `json:"tools"`
}

// InstalledPluginHookData contains the data about an installed plugin hook.
type InstalledPluginHookData struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

// InstalledPluginToolData contains the data about an installed plugin tool.
type InstalledPluginToolData struct {
	Name string `json:"name"`
}

// RegistryPlugin contains the data about a registry plugin.
type RegistryPlugin struct {
	plugin
	URL string `json:"url"`
}

// GetPlugins returns all available plugins, optionally filtered by a search query.
func (client *pluginsClient) GetPlugins() (*PluginsResponse, error) {
	resp, err := doAPIRequest[any, PluginsResponse](
		client.config,
		http.MethodGet,
		"",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UploadPluginResponse contains the information about the uploaded plugin.
type UploadPluginResponse struct {
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	Info        string `json:"info"`
}

// UploadPlugin uploads a plugin.
func (client *pluginsClient) UploadPlugin(zipFileReader *os.File) (*UploadPluginResponse, error) {
	if zipFileReader == nil {
		return nil, ErrUploadMissingFile
	}

	var requestBodyBuffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBodyBuffer)

	formFieldWriter, err := multipartWriter.CreateFormFile("file", zipFileReader.Name())
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(formFieldWriter, zipFileReader)
	if err != nil {
		return nil, err
	}

	multipartWriter.Close()

	resp, err := doHTTPRequest[UploadPluginResponse](
		client.config,
		multipartWriter.FormDataContentType(),
		http.MethodPost,
		"upload",
		nil,
		&requestBodyBuffer,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type UploadPluginFromRegistryPayload struct {
	URL string `json:"url"`
}

// UploadPluginFromRegistry uploads a plugin from a registry url.
func (client *pluginsClient) UploadPluginFromRegistry(payload UploadPluginFromRegistryPayload) (*UploadPluginResponse, error) {
	resp, err := doAPIRequest[UploadPluginFromRegistryPayload, UploadPluginResponse](
		client.config,
		http.MethodPost,
		"upload/registry",
		nil,
		&payload,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type TogglePluginResponse struct {
	Info string `json:"info"`
}

// TogglePlugin enables or disables a single plugin.
func (client *pluginsClient) TogglePlugin(pluginID string) (*TogglePluginResponse, error) {
	pathParams := fmt.Sprintf("toggle/%s", pluginID)

	resp, err := doAPIRequest[any, TogglePluginResponse](
		client.config,
		http.MethodPost,
		pathParams,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type GetPluginsSettingsResponse struct {
	Settings []PluginSetting `json:"settings"`
}

// PluginSetting contains the data about a single plugin setting.
type PluginSetting struct {
	Name   string              `json:"name"`
	Value  map[string]any      `json:"value"`
	Schema PluginSettingSchema `json:"schema"`
}

// PluginSettingSchema contains the JSON schema about a single plugin setting.
type PluginSettingSchema struct {
	settingSchema
}

// GetPluginsSettings returns the settings for all plugins.
func (client *pluginsClient) GetPluginsSettings() (*GetPluginsSettingsResponse, error) {
	resp, err := doAPIRequest[any, GetPluginsSettingsResponse](
		client.config,
		http.MethodGet,
		"settings",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetPluginSettings returns the settings for a single plugin.
func (client *pluginsClient) GetPluginSettings(pluginID string) (*PluginSetting, error) {
	pathParams := fmt.Sprintf("settings/%s", pluginID)

	resp, err := doAPIRequest[any, PluginSetting](
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

// UpsertPluginSettingsValue upserts the settings for a single plugin.
func (client *pluginsClient) UpsertPluginSettingsValue(pluginID string, value map[string]any) (*PluginSetting, error) {
	pathParams := fmt.Sprintf("settings/%s", pluginID)

	resp, err := doAPIRequest[map[string]any, PluginSetting](
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

type Plugin struct {
}

// GetPluginDetail returns the details for a single installed plugin.
func (client *pluginsClient) GetPluginDetail(pluginID string) (*InstalledPlugin, error) {
	resp, err := doAPIRequest[any, InstalledPlugin](
		client.config,
		http.MethodGet,
		pluginID,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DeletePluginResponse contains the information about the deleted plugin.
type DeletePluginResponse struct {
	// The name of the deleted plugin.
	Deleted string `json:"deleted"`
}

// DeletePlugin deletes a single plugin.
func (client *pluginsClient) DeletePlugin(pluginID string) (*DeletePluginResponse, error) {
	resp, err := doAPIRequest[any, DeletePluginResponse](
		client.config,
		http.MethodDelete,
		pluginID,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
