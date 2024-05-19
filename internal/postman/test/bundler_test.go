package postman_test

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	extPostman "github.com/giuliano-macedo/go-postman-collection"
	"github.com/giuliano-macedo/omie-docs/internal/bundler"
	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/giuliano-macedo/omie-docs/internal/postman"
	"github.com/giuliano-macedo/omie-docs/test/mocks"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name string
}

func (testCase testCase) getFilePath(fileName string) string {
	return path.Join("testdata", testCase.name, fileName)
}

func (testCase testCase) parseFile(fileName string, value any) {
	file, err := os.Open(testCase.getFilePath(fileName))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err = json.NewDecoder(file).Decode(value); err != nil {
		panic(err)
	}
}

func TestBundlerCases(t *testing.T) {
	testCaseNames := []testCase{
		{"normal"},
	}

	for _, testCase := range testCaseNames {
		t.Run(testCase.name, func(t *testing.T) {
			postmanBundler := postman.NewPostmanBundler()
			fsWriter := mocks.NewFilewriterMock()
			var (
				home               core.HomePage
				pages              []core.Page
				expectedCollection extPostman.Collection
			)
			testCase.parseFile("homepage.json", &home)
			testCase.parseFile("pages.json", &pages)
			testCase.parseFile("expected_collection.json", &expectedCollection)

			err := postmanBundler.Bundle(bundler.Args{
				Pages:    pages,
				Home:     home,
				FsWriter: fsWriter,
			})
			require.NoError(t, err)

			require.True(t, fsWriter.AllFilesClosed())
			require.Equal(t, fsWriter.FileNamesWritten(), []string{core.PostmanCollectionName})

			savedCollection, ok := fsWriter.JsonFileSaved(core.PostmanCollectionName).(extPostman.Collection)
			require.True(t, ok)

			savedCollection.Info.Name = ""
			savedCollection.Info.Description.Content = ""
			savedCollection.Info.Version = ""

			require.Equal(t, expectedCollection, savedCollection)

		})
	}
}
