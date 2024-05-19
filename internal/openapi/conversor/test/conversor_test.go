package openapi_conversor_test

import (
	"encoding/json"
	"flag"
	"os"
	"path"
	"testing"

	"github.com/giuliano-macedo/omie-docs/internal/core"
	openapi_conversor "github.com/giuliano-macedo/omie-docs/internal/openapi/conversor"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/openapi-go/openapi3"
)

type testCase struct {
	name         string
	featureIndex int
}

func (testCase testCase) getFilePath(fileName string) string {
	return path.Join("test_data", testCase.name, fileName)
}

func (testCase testCase) parseFile(fileName string, value any) {
	data, err := os.ReadFile(testCase.getFilePath(fileName))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, value)
	if err != nil {
		panic(err)
	}
}

func readTestCase(testCase testCase) (homepage core.HomePage, pages []core.Page, feature *core.Feature, expectedSpec interface{}) {
	testCase.parseFile("homepage.json", &homepage)
	testCase.parseFile("pages.json", &pages)
	testCase.parseFile("expected_openapi.json", &expectedSpec)

	if testCase.featureIndex >= 0 {
		feature = &homepage.Features[testCase.featureIndex]
	}
	return
}

func openapiSpecToInterface(spec openapi3.Spec) (specInterface interface{}) {
	data, err := spec.MarshalJSON()
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &specInterface)
	if err != nil {
		panic(err)
	}
	return
}

var updateCases = flag.Bool("update_cases", false, "Update expected openapi cases")

func TestConversorCases(t *testing.T) {
	flag.Parse()
	testCases := []testCase{
		{"all_types", -1},
		{"no_response", 0},
		{"required_and_deprecated_fields", 0},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			homepage, pages, feature, expectedSpec := readTestCase(testCase)
			spec := openapi_conversor.ConvertToOpenApi(homepage, pages, feature)
			spec.Info.Description = nil
			spec.Info.Version = "0"
			if *updateCases {
				data, _ := json.MarshalIndent(openapiSpecToInterface(spec), "", "  ")
				err := os.WriteFile(testCase.getFilePath("expected_openapi.json"), data, os.ModePerm)
				if err != nil {
					panic(err)
				}
				return
			}
			if !assert.Equal(t, expectedSpec, openapiSpecToInterface(spec)) {
				t.Error("spec not matching, run with -update_cases=true if you are sure it is generating the correct spec")
			}
		})
	}
}
