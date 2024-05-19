package bruno_test

import (
	"testing"

	extBruno "github.com/giuliano-macedo/go-bruno-collection/pkg/bruno"
	"github.com/giuliano-macedo/omie-docs/internal/bruno"
	"github.com/giuliano-macedo/omie-docs/internal/bruno/test/testdata"
	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/stretchr/testify/require"
)

func TestConversorCases(t *testing.T) {
	testCaseNames := []string{
		"normal",
	}

	for _, testCase := range testCaseNames {
		t.Run(testCase, func(t *testing.T) {
			home := testdata.ReadJsonFile[core.HomePage](testCase, "homepage.json")
			pages := testdata.ReadJsonFile[[]core.Page](testCase, "pages.json")
			expectedCollection := testdata.ReadJsonFile[extBruno.Collection](testCase, "expected_collection.json")

			savedCollection := bruno.ConvertToBruno(home, pages)
			savedCollection.Name = ""
			savedCollection.Docs = ""

			require.Equal(t, expectedCollection, savedCollection)
		})
	}
}
