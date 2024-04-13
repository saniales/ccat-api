package ccatapi

import (
	"fmt"
	"net/http"
	"net/url"
)

// memoryClient is a sub-client for the Memory API.
type memoryClient struct {
	config clientConfig
}

// newMemoryClient creates a new Memory sub-client with the provided config.
func newMemoryClient(config clientConfig) *memoryClient {
	client := &memoryClient{
		config: config,
	}

	WithBaseURL(fmt.Sprintf("%s/%s", client.config.baseURL, "memory"))(&client.config)

	return client
}

type RecallMemoriesResponse struct {
	Query   recallMemoriesResponseQuery   `json:"query"`
	Vectors recallMemoriesResponseVectors `json:"vectors"`
}

type recallMemoriesResponseQuery struct {
	Text   string    `json:"text"`
	Vector []float64 `json:"vector"`
}

type recallMemoriesResponseVectors struct {
	Embedder    string                                   `json:"embedder"`
	Collections recallMemoriesResponseVectorsCollections `json:"collections"`
}

type recallMemoriesResponseVectorsCollections struct {
	Episodic    []Memory `json:"episodic"`
	Declarative []Memory `json:"declarative"`
	Procedural  []Memory `json:"procedural"`
}

type Memory struct {
	PageContent string         `json:"page_content"`
	Metadata    memoryMetadata `json:"metadata"`
	Type        string         `json:"type"`
	ID          string         `json:"id"`
	Score       float64        `json:"score"`
	Vector      []float64      `json:"vector"`
}

type memoryMetadata struct {
	Source string  `json:"source"`
	When   float64 `json:"when"`
}

// RecallMemories searches memories similar to given text.
func (client *memoryClient) RecallMemories(text string, k uint) (*RecallMemoriesResponse, error) {
	queryParams := make(url.Values, 2)

	queryParams.Set("text", text)
	queryParams.Set("k", fmt.Sprint(k))

	resp, err := doAPIRequest[any, RecallMemoriesResponse](
		client.config,
		http.MethodGet,
		"recall",
		queryParams,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetMemoryCollectionsResponse contains the response of a GetMemoryCollections call.
type GetMemoryCollectionsResponse struct {
	// The available collections
	Collections []MemoryCollection `json:"collections"`
}

// MemoryCollection contains the data about a single memory collection.
type MemoryCollection struct {
	// The name of the collection
	Name string `json:"name"`

	// The number of vectors in the collection
	VectorsCount uint `json:"vectors_count"`
}

// GetMemoryCollections returns all memories collections data.
func (client *memoryClient) GetMemoryCollections() (*GetMemoryCollectionsResponse, error) {
	resp, err := doAPIRequest[any, GetMemoryCollectionsResponse](
		client.config,
		http.MethodGet,
		"collections",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// WipeMemoryCollectionsResponse contains the response of a WipeMemoryCollections call.
type WipeMemoryCollectionsResponse struct {
	Episodic    bool `json:"episodic,omitempty"`
	Declarative bool `json:"declarative,omitempty"`
	Procedural  bool `json:"procedural,omitempty"`
}

// WipeMemoryCollections wipes all memories collections data.
func (client *memoryClient) WipeMemoryCollections() (*WipeMemoryCollectionsResponse, error) {
	resp, err := doAPIRequest[any, WipeMemoryCollectionsResponse](
		client.config,
		http.MethodDelete,
		"collections",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// WipeMemoryCollection wipes all memories in a collection.
func (client *memoryClient) WipeMemoryCollection(id string) (*WipeMemoryCollectionsResponse, error) {
	pathParams := fmt.Sprintf("collections/%s", id)

	resp, err := doAPIRequest[any, WipeMemoryCollectionsResponse](
		client.config,
		http.MethodDelete,
		pathParams,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// WipeMemoryCollectionPointResponse contains the response of a WipeMemoryCollectionPoint call.
type WipeMemoryCollectionPointResponse struct {
	// The ID of the collected point
	Deleted string `json:"deleted"`
}

// WipeMemoryCollectionPoint wipes a single memory in a collection.
func (client *memoryClient) WipeMemoryCollectionPoint(collectionID string, memoryID string) (*WipeMemoryCollectionsResponse, error) {
	pathParams := fmt.Sprintf("collections/%s/points/%s", collectionID, memoryID)

	resp, err := doAPIRequest[any, WipeMemoryCollectionsResponse](
		client.config,
		http.MethodDelete,
		pathParams,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// WipeMemoryCollectionPointsByMetadata wipes all memories in a collection by metadata.
func (client *memoryClient) WipeMemoryCollectionPointsByMetadata(collectionID string, metadata map[string]any) (*WipeMemoryCollectionsResponse, error) {
	pathParams := fmt.Sprintf("collections/%s/points", collectionID)

	resp, err := doAPIRequest[map[string]any, WipeMemoryCollectionsResponse](
		client.config,
		http.MethodDelete,
		pathParams,
		nil,
		&metadata,
	)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

type GetConversationHistoryResponse struct {
	History []conversationMessage `json:"history"`
}

type conversationMessage struct {
	Who     string                 `json:"who"`
	Message string                 `json:"message"`
	Why     conversationMessageWhy `json:"why"`
}

type conversationMessageWhy struct {
	Input             string                          `json:"input"`
	IntermediateSteps []any                           `json:"intermediate_steps"`
	Memory            []conversationMessageMemoryData `json:"memory"`
}

type conversationMessageMemoryData struct {
	Episodic    []any `json:"episodic"`
	Declarative []any `json:"declarative"`
	Procedural  []any `json:"procedural"`
}

// GetConversationHistory gets all conversation histories.
func (client *memoryClient) GetConversationHistory() (*GetConversationHistoryResponse, error) {
	resp, err := doAPIRequest[any, GetConversationHistoryResponse](
		client.config,
		http.MethodGet,
		"conversation_history",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// WipeConversationHistoryResponse contains the response of a WipeConversationHistory call.
type WipeConversationHistoryResponse struct {
	Deleted bool `json:"deleted"`
}

// WipeConversationHistory wipes all conversation history.
func (client *memoryClient) WipeConversationHistory() (*WipeConversationHistoryResponse, error) {
	resp, err := doAPIRequest[any, WipeConversationHistoryResponse](
		client.config,
		http.MethodDelete,
		"conversation_history",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
