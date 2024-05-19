package crawler

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"testing"

	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/giuliano-macedo/omie-docs/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type pageTestCase struct {
	name string
	url  string
}

func TestExtractPageUrlInfoErrors(t *testing.T) {
	page := new(core.Page)

	err := extractPageUrlInfo("notanurl.\\:not an url", page)
	assert.Error(t, err)

	err = extractPageUrlInfo("https://example.com/justonecomponent", page)
	assert.Error(t, err)
}

func TestGetPageErrorOnFetcher(t *testing.T) {
	fetcherErr := errors.New("fetcher error")
	var fetcherMock = new(mocks.FetcherMock)
	fetcherMock.On("GetPage", mock.Anything).Return(io.ReadCloser(nil), fetcherErr)

	_, err := GetPage(fetcherMock, "")
	if !assert.Error(t, err) {
		assert.Equal(t, fetcherErr, err)
	}
}

func TestGetPageCases(t *testing.T) {
	cases := []pageTestCase{
		{"all_types_case", "https://example.com/base/version1/directory/to/file/"},
		{"missing_fields_case", ""},
		{"deprecated_case", ""},
		{"non_object_array_case", ""},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			runPageTestCase(t, testCase)
		})
	}
}

func TestGetPageErrCases(t *testing.T) {
	cases := []pageTestCase{
		{"no_page_header_err_case", ""},
		{"multiple_params_err_case", ""},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			runTestCaseErr(t, testCase)
		})
	}
}

func getPageTestSuiteData(caseName string) (expectedPage core.Page, pageData string, err error) {
	data, err := os.ReadFile(path.Join("test_data/page", caseName, "page.html"))
	if err != nil {
		return
	}

	pageData = string(data)
	expectedPageData, err := os.ReadFile(path.Join("test_data/page", caseName, "expected_page.json"))
	if err != nil {
		return
	}

	err = json.Unmarshal(expectedPageData, &expectedPage)

	return
}

func runPageTestCase(t *testing.T, testCase pageTestCase) {
	expectedPage, pageData, err := getPageTestSuiteData(testCase.name)
	if !assert.NoError(t, err) {
		return
	}

	var fetcherMock = new(mocks.FetcherMock)
	fetcherMock.On("GetPage", mock.Anything).Return(mocks.NewReadCloseMock(string(pageData)), nil)
	page, err := GetPage(fetcherMock, testCase.url)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, expectedPage, page)
}

func runTestCaseErr(t *testing.T, testCase pageTestCase) {
	pageData, err := os.ReadFile(path.Join("test_data/page", testCase.name, "page.html"))
	if !assert.NoError(t, err) {
		return
	}

	var fetcherMock = new(mocks.FetcherMock)
	fetcherMock.On("GetPage", mock.Anything).Return(mocks.NewReadCloseMock(string(pageData)), nil)

	_, err = GetPage(fetcherMock, testCase.url)
	if !assert.Error(t, err) {
		return
	}
}
