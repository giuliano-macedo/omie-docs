package openapi_bundler_test

import (
	"errors"
	"testing"

	"github.com/giuliano-macedo/omie-docs/internal/bundler"
	"github.com/giuliano-macedo/omie-docs/internal/core"
	openapi_bundler "github.com/giuliano-macedo/omie-docs/internal/openapi/bundler"
	"github.com/giuliano-macedo/omie-docs/test/mocks"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/openapi-go/openapi3"
)

func getPagesAndHome() ([]core.Page, core.HomePage) {
	pages := []core.Page{
		{Name: "page1", FullUrl: "example.com/page1"},
	}
	home := core.HomePage{
		Features: []core.Feature{
			{
				Name: "Feature 1",
				MainEntities: []core.Entity{
					{Url: "example.com/page1"},
				},
			},
		},
	}
	return pages, home
}

func TestBundler(t *testing.T) {
	fsWriter := mocks.NewFilewriterMock()
	pages, home := getPagesAndHome()
	openApiBundler := openapi_bundler.NewOpenApiBundler()
	err := openApiBundler.Bundle(bundler.Args{
		Pages:    pages,
		Home:     home,
		FsWriter: fsWriter,
	})

	require.NoError(t, err)

	require.True(t, fsWriter.AllFilesClosed())
	require.Equal(t, []string{"openapi"}, fsWriter.DirectoriesCreated())
	require.Equal(t, []string{"index.html", "openapi/all.openapi.json", "openapi/feature_1.openapi.json"}, fsWriter.FileNamesWritten())

	requireToBeValidOpenApi(t, fsWriter.JsonFileSaved("openapi/feature_1.openapi.json"))
	requireToBeValidOpenApi(t, fsWriter.JsonFileSaved("openapi/all.openapi.json"))

	indexContent := string(fsWriter.FileWritten("index.html").FileContent)

	require.Contains(t, indexContent, "openapi/feature_1.openapi.json")
	require.Contains(t, indexContent, "openapi/all.openapi.json")

}

func TestBundlerErrors(t *testing.T) {
	runTest := func(name string, configFileWritterMock func(fsWriter *mocks.FileWriterMock, expectedError error)) {
		t.Run(name, func(t *testing.T) {
			expectedError := errors.New(name)
			pages, home := getPagesAndHome()
			fsWriter := mocks.NewFilewriterMock()
			configFileWritterMock(fsWriter, expectedError)

			openApiBundler := openapi_bundler.NewOpenApiBundler()
			err := openApiBundler.Bundle(bundler.Args{
				Pages:    pages,
				Home:     home,
				FsWriter: fsWriter,
			})

			require.ErrorIs(t, err, expectedError)
		})
	}

	runTest("MkDirAll", func(fsWriter *mocks.FileWriterMock, expectedError error) {
		fsWriter.ErrorOnMkdirAll("openapi", expectedError)
	})

	runTest("Save OpenApi", func(fsWriter *mocks.FileWriterMock, expectedError error) {
		fsWriter.ErrorOnSaveJson("openapi/feature_1.openapi.json", expectedError)
	})

	runTest("Save OpenApi All", func(fsWriter *mocks.FileWriterMock, expectedError error) {
		fsWriter.ErrorOnSaveJson("openapi/all.openapi.json", expectedError)
	})

	runTest("Create index", func(fsWriter *mocks.FileWriterMock, expectedError error) {
		fsWriter.ErrorOnCreate("index.html", expectedError)
	})

	runTest("Write index", func(fsWriter *mocks.FileWriterMock, expectedError error) {
		fsWriter.ErrorOnWrite("index.html", expectedError)
	})
}

func requireToBeValidOpenApi(t *testing.T, value interface{}) {
	_, ok := value.(openapi3.Spec)
	require.True(t, ok)
}
