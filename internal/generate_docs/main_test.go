//go:build integration
// +build integration

package generate_docs_test

import (
	"os"
	"testing"

	"github.com/giuliano-macedo/omie-docs/internal/generate_docs"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/openapi-go/openapi3"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func readOpenapiJsonFile(t *testing.T, path string) (spec openapi3.Spec) {
	data, err := os.ReadFile(path)
	assert.NoError(t, err)

	err = spec.UnmarshalJSON(data)
	assert.NoError(t, err)

	return
}

func TestMain(t *testing.T) {
	vcr, err := recorder.NewWithOptions(
		&recorder.Options{
			CassetteName: "../../integration_test_data/fixtures/omie",
			Mode:         recorder.ModeReplayWithNewEpisodes,
		},
	)
	assert.NoError(t, err)
	defer vcr.Stop()

	generate_docs.Run(generate_docs.Args{
		HttpClient:                vcr.GetDefaultClient(),
		OutputDir:                 "../../integration_test_data/bundle",
		NumberOfWorkers:           16,
		DumpIntermediaryDataTypes: false,
		CachingDirectory:          "",
		InitialPage:               "",
		PrettifyJson:              true,
	})
	expectedOpenapiSchema := readOpenapiJsonFile(t, "../../integration_test_data/expected_all.openapi.json")
	actualOpenApiSchema := readOpenapiJsonFile(t, "../../integration_test_data/bundle/openapi/all.openapi.json")

	// check not needed
	actualOpenApiSchema.Info.Description = nil
	expectedOpenapiSchema.Info.Description = nil

	// error message too long when fields missmatch
	assert.Equal(t, expectedOpenapiSchema, actualOpenApiSchema)
}
