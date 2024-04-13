package ccatapi

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// rabbitHoleClient is a sub-client for the Rabbit Hole API.
type rabbitHoleClient struct {
	config clientConfig
}

// newRabbitHoleClient creates a new Rabbit Hole sub-client with the provided config.
func newRabbitHoleClient(config clientConfig) *rabbitHoleClient {
	client := &rabbitHoleClient{
		config: config,
	}

	WithBaseURL(fmt.Sprintf("%s/%s", client.config.baseURL, "rabbit_hole"))(&client.config)

	return client
}

// UploadPayload is the payload for the upload endpoint.
type UploadPayload struct {
	File         *os.File `json:"file"`
	ChunkSize    int      `json:"chunk_size"`
	ChunkOverlap int      `json:"chunk_overlap"`
}

// UploadResponse is the response for the upload endpoint.
type UploadResponse struct {
	URL  string `json:"url"`
	Info string `json:"info"`
}

// Upload uploads a file into the rabbit hole.
func (client *rabbitHoleClient) Upload(payload UploadPayload) (*UploadResponse, error) {
	if payload.File == nil {
		return nil, ErrUploadMissingFile
	}

	var requestBodyBuffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBodyBuffer)

	formFieldWriter, err := multipartWriter.CreateFormFile("file", payload.File.Name())
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(formFieldWriter, payload.File)
	if err != nil {
		return nil, err
	}

	err = multipartWriter.WriteField("chunk_size", fmt.Sprint(payload.ChunkSize))
	if err != nil {
		return nil, err
	}

	err = multipartWriter.WriteField("chunk_overlap", fmt.Sprint(payload.ChunkOverlap))
	if err != nil {
		return nil, err
	}

	multipartWriter.Close()

	resp, err := doHTTPRequest[UploadResponse](
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

// UploadFromURLPayload is the payload for the UploadFromURL endpoint.
type UploadFromURLPayload struct {
	URL          string `json:"url"`
	ChunkSize    int    `json:"chunk_size"`
	ChunkOverlap int    `json:"chunk_overlap"`
}

// UploadFromURLResponse is the response for the UploadFromURL endpoint.
type UploadFromURLResponse struct {
	URL  string `json:"url"`
	Info string `json:"info"`
}

// UploadFromURL uploads a file from a URL into the rabbit hole.
func (client *rabbitHoleClient) UploadFromURL(payload UploadFromURLPayload) (*UploadFromURLResponse, error) {
	resp, err := doAPIRequest[UploadFromURLPayload, UploadFromURLResponse](
		client.config,
		http.MethodPost,
		"upload",
		nil,
		&payload,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UploadMemoryResponse represents the response from the UploadMemory call.
type UploadMemoryResponse struct {
	File string `json:"file"`
}

// UploadMemory uploads a memory into the rabbit hole.
func (client *rabbitHoleClient) UploadMemory(memoryFile *os.File) (*UploadMemoryResponse, error) {
	if memoryFile == nil {
		return nil, ErrUploadMissingFile
	}

	var requestBodyBuffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBodyBuffer)

	formFieldWriter, err := multipartWriter.CreateFormFile("file", memoryFile.Name())
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(formFieldWriter, memoryFile)
	if err != nil {
		return nil, err
	}

	multipartWriter.Close()

	resp, err := doHTTPRequest[UploadMemoryResponse](
		client.config,
		multipartWriter.FormDataContentType(),
		http.MethodPost,
		"memory",
		nil,
		&requestBodyBuffer,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetAllowedMIMETypesResponse represents the response from the GetAllowedMIMETypes call.
type GetAllowedMIMETypesResponse struct {
	Allowed []string `json:"allowed"`
}

// GetAllowedMIMETypes returns the list of allowed mime types for upload into the rabbit hole.
func (client *rabbitHoleClient) GetAllowedMIMETypes() (*GetAllowedMIMETypesResponse, error) {
	resp, err := doAPIRequest[any, GetAllowedMIMETypesResponse](
		client.config,
		http.MethodGet,
		"allowed-mimetypes",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
