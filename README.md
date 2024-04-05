# Cheshire Cat API Golang Client

This is a very simple Golang client to call [Cheshire Cat API](https://cheshirecat.ai).

The library is documented using [GoDoc](https://pkg.go.dev/github.com/saniales/ccat-api).

## Download and usage

``` bash
go get github.com/saniales/ccat-api
```

## A more complete usage example

``` go
package ccatapi_test

import (
	"fmt"
	"log"

	ccatapi "github.com/saniales/ccat-api"
)

func main() {
	// Create a new Cheshire Cat API client.
	client := ccatapi.NewClient(
		ccatapi.WithBaseURL("https://examplecat.ai"),
	)

	// Call the Cheshire Cat API
	err := client.Status()
	if err != nil {
		log.Fatal("Cheshire API is not OK")
	}

	// Use Settings API
	getSettingsParams := ccatapi.GetSettingsParams{
		Search: "example",
	}
	getSettingsResponse, err := client.Settings.GetSettings(getSettingsParams)
	if err != nil {
		log.Fatal("Cannot get settings", err)
	}
	fmt.Println(getSettingsResponse.Settings)

	// Use LLMs API
	getAllLLMsSettingsResponse, err := client.LLMs.GetAllLLMsSettings()
	if err != nil {
		log.Fatal("Cannot get LLMs settings", err)
	}
	fmt.Println(getAllLLMsSettingsResponse.Settings)

	// Use Embedders API
	getAllEmbeddersSettingsResponse, err := client.Embedders.GetAllEmbeddersSettings()
	if err != nil {
		log.Fatal("Cannot get Embedders settings", err)
	}
	fmt.Println(getAllEmbeddersSettingsResponse.Settings)

    // and so on...
}
```
