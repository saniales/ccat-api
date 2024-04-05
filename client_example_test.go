package ccatapi_test

import (
	"fmt"
	"net/http"

	ccatapi "github.com/saniales/ccat-api"
)

func ExampleNewClient_defaults() {
	// Create a new Cheshire Cat API client with only the default values.
	client := ccatapi.NewClient()

	// Call the Cheshire Cat API
	fmt.Println(client.Status())
}

func ExampleNewClient_options() {
	// Create a new Cheshire Cat API client by changing the defaults using
	// options.
	client := ccatapi.NewClient(
		ccatapi.WithHTTPClient(&http.Client{}),
		ccatapi.WithBaseURL("https://localhost:1865"),
		ccatapi.WithUserAgent("cheshire-gopher-api"),
		ccatapi.WithUserID("my_user"),
	)

	// Call the Cheshire Cat API with new settings
	fmt.Println(client.Status())
}

func ExampleClient_Status() {
	// Create a new Cheshire Cat API client.
	client := ccatapi.NewClient()

	// Call the Cheshire Cat API
	fmt.Println(client.Status())
}
